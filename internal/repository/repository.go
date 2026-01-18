package repository

import (
	"indexer/internal/model"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	// Blocks & Transactions
	GetLatestBlocks(chain model.ChainType, limit, offset int) ([]model.Block, error)
	GetBlockByHeight(chain model.ChainType, height uint64) (*model.Block, error)
	GetTransactions(chain model.ChainType, limit, offset int) ([]model.Transaction, error)
	CountTransactions(chain model.ChainType) (int64, error)

	// Sync Logic
	GetState(chain model.ChainType) (uint64, error)
	GetOrCreateState(chain model.ChainType, latestBlock uint64, configuredStart int) (uint64, error)
	SaveBlockWithTransactions(block *model.Block, txs []*model.Transaction) error

	// Read Logic (New)
	CountBlocks(chain model.ChainType) (int64, error)
	GetTransactionsByBlock(chain model.ChainType, height uint64) ([]model.Transaction, error)
	FindTransactionByHash(chain model.ChainType, hash string) (*model.Transaction, error)
	GetMaxBlockHeight(chain model.ChainType) (uint64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Table helpers
func (r *repository) blockTable(chain model.ChainType) string {
	if chain == model.ChainBTC {
		return "btc_blocks"
	}
	return "eth_blocks"
}

func (r *repository) txTable(chain model.ChainType) string {
	if chain == model.ChainBTC {
		return "btc_transactions"
	}
	return "eth_transactions"
}

// API READ METHODS
func (r *repository) GetLatestBlocks(chain model.ChainType, limit, offset int) ([]model.Block, error) {
	var blocks []model.Block
	err := r.db.Table(r.blockTable(chain)).
		Order("height DESC").
		Limit(limit).
		Offset(offset).
		Find(&blocks).Error
	return blocks, err
}

func (r *repository) GetBlockByHeight(chain model.ChainType, height uint64) (*model.Block, error) {
	var block model.Block
	err := r.db.Table(r.blockTable(chain)).Where("height = ?", height).First(&block).Error
	if err != nil {
		return nil, err
	}
	return &block, nil
}

func (r *repository) GetTransactions(chain model.ChainType, limit, offset int) ([]model.Transaction, error) {
	var txs []model.Transaction
	err := r.db.Table(r.txTable(chain)).
		Order("block_height DESC").
		Limit(limit).
		Offset(offset).
		Find(&txs).Error
	return txs, err
}

func (r *repository) CountTransactions(chain model.ChainType) (int64, error) {
	var count int64
	err := r.db.Table(r.txTable(chain)).Count(&count).Error
	return count, err
}

// SYNC METHODS
func (r *repository) GetState(chain model.ChainType) (uint64, error) {
	var state model.IndexerState
	err := r.db.Where("chain = ?", chain).Limit(1).Find(&state).Error
	if err != nil {
		return 0, err
	}
	return state.LastIndexedHeight, nil
}

func (r *repository) GetOrCreateState(chain model.ChainType, latestBlock uint64, configuredStart int) (uint64, error) {
	var state model.IndexerState
	err := r.db.Where("chain = ?", chain).Limit(1).Find(&state).Error
	if err != nil {
		return 0, err
	}

	// 1. If record exists, we might still need to reset if it's invalid (e.g. height > tip)
	if state.Chain != "" {
		// User requirement: If height is greater than current chain tip → reset to (latest - 10)
		if state.LastIndexedHeight > latestBlock {
			resetHeight := uint64(0)
			if latestBlock > 10 {
				resetHeight = latestBlock - 10
			}
			state.LastIndexedHeight = resetHeight
			state.UpdatedAt = time.Now()
			r.db.Save(&state)
			log.Printf("[%s] Reset state to %d because last indexed height was greater than tip", chain, resetHeight)
		}
		return state.LastIndexedHeight, nil
	}

	// 2. Not found, create it
	startHeight := uint64(0)

	// User requirement: If configured start height is zero or negative → start from (latest - 10)
	if configuredStart <= 0 {
		if latestBlock > 10 {
			startHeight = latestBlock - 10
		}
	} else {
		startHeight = uint64(configuredStart)
		// User requirement: If configured start height is greater than tip → reset to (latest - 10)
		if startHeight > latestBlock {
			if latestBlock > 10 {
				startHeight = latestBlock - 10
			} else {
				startHeight = 0
			}
		}
	}

	state = model.IndexerState{
		Chain:             chain,
		LastIndexedHeight: startHeight,
		UpdatedAt:         time.Now(),
	}

	if err := r.db.Create(&state).Error; err != nil {
		return 0, err
	}

	log.Printf("[%s] Created initial indexer state starting at block %d", chain, startHeight)
	return startHeight, nil
}

func (r *repository) SaveBlockWithTransactions(block *model.Block, txs []*model.Transaction) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Save Block
		if err := tx.Table(r.blockTable(block.Chain)).Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(block).Error; err != nil {
			return err
		}

		// 2. Save Transactions
		if len(txs) > 0 {
			if err := tx.Table(r.txTable(block.Chain)).Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(txs).Error; err != nil {
				return err
			}
		}

		// 3. Update State
		if err := tx.Model(&model.IndexerState{}).
			Where("chain = ?", block.Chain).
			Update("last_indexed_block", block.Height).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *repository) CountBlocks(chain model.ChainType) (int64, error) {
	var count int64
	err := r.db.Table(r.blockTable(chain)).Count(&count).Error
	return count, err
}

func (r *repository) GetTransactionsByBlock(chain model.ChainType, height uint64) ([]model.Transaction, error) {
	var txs []model.Transaction
	err := r.db.Table(r.txTable(chain)).Where("block_height = ?", height).Find(&txs).Error
	return txs, err
}

func (r *repository) FindTransactionByHash(chain model.ChainType, hash string) (*model.Transaction, error) {
	var tx model.Transaction
	err := r.db.Table(r.txTable(chain)).Where("hash = ?", hash).First(&tx).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *repository) GetMaxBlockHeight(chain model.ChainType) (uint64, error) {
	var max uint64
	// Use Scan to handle potential NULL result if table is empty
	err := r.db.Table(r.blockTable(chain)).Select("COALESCE(MAX(height), 0)").Row().Scan(&max)
	return max, err
}
