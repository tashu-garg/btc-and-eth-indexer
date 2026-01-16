package btc

import (
	"time"

	"indexer/internal/constants"
	"indexer/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	GetLastIndexedBlock() (uint64, error)
	UpdateIndexerState(blockHeight uint64) error
	SaveBlock(block *model.Block) error
	SaveTransaction(tx *model.Transaction) error
	GetBlocks(limit, offset int) ([]model.Block, error)
	GetBlockByHeight(height uint64) (*model.Block, error)
	GetBlockByHash(hash string) (*model.Block, error)
	GetTransactions(limit, offset int) ([]model.Transaction, error)
	GetTransactionByHash(hash string) (*model.Transaction, error)
	DeleteBlocksFrom(height uint64) error
	InitIndexerState(chain string) error
	GetStats() (map[string]interface{}, error)
	SearchBlocks(query string) ([]model.Block, error)
	SearchTransactions(query string) ([]model.Transaction, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// ================= INDEXER STATE =================

func (r *repository) GetLastIndexedBlock() (uint64, error) {
	var state model.IndexerState
	// Use Limit(1).Find to avoid "record not found" noise in logs
	err := r.db.Where("chain = ?", constants.ChainBTC).Limit(1).Find(&state).Error
	if err != nil {
		return 0, err
	}
	if state.Chain == "" {
		return 0, nil // no record yet, start from 0
	}
	return state.LastIndexedHeight, nil
}

func (r *repository) UpdateIndexerState(blockHeight uint64) error {
	now := time.Now()
	return r.db.Exec(`
		INSERT INTO indexer_states (chain, last_indexed_block, updated_at)
		VALUES (?, ?, ?)
		ON CONFLICT (chain)
		DO UPDATE SET last_indexed_block = EXCLUDED.last_indexed_block,
		              updated_at = EXCLUDED.updated_at
	`, constants.ChainBTC, blockHeight, now).Error
}

func (r *repository) InitIndexerState(chain string) error {
	state := model.IndexerState{}
	err := r.db.Where("chain = ?", chain).First(&state).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			state = model.IndexerState{
				Chain:             model.ChainType(chain),
				LastIndexedHeight: 0,
				UpdatedAt:         time.Now(),
			}
			return r.db.Create(&state).Error
		}
		return err
	}
	return nil
}

func (r *repository) SaveBlock(block *model.Block) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "chain"}, {Name: "height"}},
		DoNothing: true,
	}).Create(block).Error
}

func (r *repository) GetBlocks(limit, offset int) ([]model.Block, error) {
	var blocks []model.Block
	err := r.db.Where("chain = ?", constants.ChainBTC).
		Order("height DESC").
		Limit(limit).
		Offset(offset).
		Find(&blocks).Error
	return blocks, err
}

func (r *repository) GetBlockByHeight(height uint64) (*model.Block, error) {
	var block model.Block
	err := r.db.Where("chain = ? AND height = ?", constants.ChainBTC, height).First(&block).Error
	if err != nil {
		return nil, err
	}
	return &block, nil
}

func (r *repository) GetBlockByHash(hash string) (*model.Block, error) {
	var block model.Block
	err := r.db.Preload("Transactions").Where("chain = ? AND hash = ?", constants.ChainBTC, hash).First(&block).Error
	if err != nil {
		return nil, err
	}
	return &block, nil
}

// ================= TRANSACTIONS =================

func (r *repository) SaveTransaction(tx *model.Transaction) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "hash"}},
		DoNothing: true,
	}).Create(tx).Error
}

func (r *repository) GetTransactions(limit, offset int) ([]model.Transaction, error) {
	var txs []model.Transaction
	err := r.db.Where("chain = ?", constants.ChainBTC).
		Order("block_height DESC").
		Limit(limit).
		Offset(offset).
		Find(&txs).Error
	return txs, err
}

func (r *repository) GetTransactionByHash(hash string) (*model.Transaction, error) {
	var tx model.Transaction
	err := r.db.Where("chain = ? AND hash = ?", constants.ChainBTC, hash).First(&tx).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

// ================= REORG HANDLING =================

func (r *repository) DeleteBlocksFrom(height uint64) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Exec(`DELETE FROM blocks WHERE chain = ? AND height >= ?`, constants.ChainBTC, height).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Exec(`DELETE FROM transactions WHERE chain = ? AND block_height >= ?`, constants.ChainBTC, height).Error; err != nil {
		tx.Rollback()
		return err
	}

	var newLast uint64 = 0
	if height > 0 {
		newLast = height - 1
	}
	if err := tx.Exec(`UPDATE indexer_states SET last_indexed_block = ? WHERE chain = ?`, newLast, constants.ChainBTC).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
func (r *repository) GetStats() (map[string]interface{}, error) {
	var blockCount int64
	var txCount int64

	if err := r.db.Model(&model.Block{}).Where("chain = ?", constants.ChainBTC).Count(&blockCount).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&model.Transaction{}).Where("chain = ?", constants.ChainBTC).Count(&txCount).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_blocks":       blockCount,
		"total_transactions": txCount,
	}, nil
}

func (r *repository) SearchBlocks(query string) ([]model.Block, error) {
	var blocks []model.Block
	err := r.db.Where("chain = ? AND (hash ILIKE ? OR CAST(height AS TEXT) LIKE ?)",
		constants.ChainBTC, "%"+query+"%", query+"%").
		Order("height DESC").Limit(10).Find(&blocks).Error
	return blocks, err
}

func (r *repository) SearchTransactions(query string) ([]model.Transaction, error) {
	var txs []model.Transaction
	// BTC uses from_address / to_address columns
	err := r.db.Where("chain = ? AND (hash ILIKE ? OR from_address ILIKE ? OR to_address ILIKE ?)",
		constants.ChainBTC, "%"+query+"%", "%"+query+"%", "%"+query+"%").
		Order("block_height DESC").Limit(10).Find(&txs).Error
	return txs, err
}
