package eth

import (
	"errors"
	"indexer/internal/constants"
	"indexer/internal/model"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	GetLastIndexedBlock() (uint64, error)
	UpdateIndexerState(blockHeight uint64) error
	SaveBlock(block *model.Block) error
	SaveTransactions(txs []model.Transaction) error
	GetBlockByHeight(height uint64) (*model.Block, error)
	GetBlockByHash(hash string) (*model.Block, error)
	GetBlocks(limit, offset int) ([]model.Block, error)
	GetTransactions(limit, offset int) ([]model.Transaction, error)
	GetTransactionByHash(hash string) (*model.Transaction, error)
	DeleteBlocksFrom(height uint64) error
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
	err := r.db.Where("chain = ?", constants.ChainETH).Limit(1).Find(&state).Error
	if err != nil {
		return 0, err
	}

	// If ID is 0 or Chain is empty, it means no record was found
	if state.Chain == "" {
		// Initialize state if not found
		state = model.IndexerState{
			Chain:             constants.ChainETH,
			LastIndexedHeight: 10056000,
			UpdatedAt:         time.Now(),
		}
		if err := r.db.Create(&state).Error; err != nil {
			return 10056000, err
		}
		return 10056000, nil
	}
	return state.LastIndexedHeight, nil
}

func (r *repository) UpdateIndexerState(blockHeight uint64) error {
	return r.db.Model(&model.IndexerState{}).
		Where("chain = ?", constants.ChainETH).
		Updates(map[string]interface{}{
			"last_indexed_block": blockHeight,
			"updated_at":         time.Now(),
		}).Error
}

// ================= BLOCKS =================

func (r *repository) SaveBlock(block *model.Block) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "chain"}, {Name: "height"}},
		DoNothing: true,
	}).Create(block).Error
}

func (r *repository) GetBlockByHeight(height uint64) (*model.Block, error) {
	var block model.Block
	err := r.db.Where("chain = ? AND height = ?", constants.ChainETH, height).First(&block).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &block, nil
}

func (r *repository) GetBlocks(limit, offset int) ([]model.Block, error) {
	var blocks []model.Block
	err := r.db.Where("chain = ?", constants.ChainETH).
		Order("height DESC").
		Limit(limit).
		Offset(offset).
		Find(&blocks).Error
	return blocks, err
}

// ================= TRANSACTIONS =================

func (r *repository) SaveTransactions(txs []model.Transaction) error {
	if len(txs) == 0 {
		return nil
	}
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "hash"}},
		DoNothing: true,
	}).Create(&txs).Error
}

func (r *repository) GetTransactions(limit, offset int) ([]model.Transaction, error) {
	var txs []model.Transaction
	err := r.db.Where("chain = ?", constants.ChainETH).
		Order("block_height DESC").
		Limit(limit).
		Offset(offset).
		Find(&txs).Error
	return txs, err
}

// ================= REORG HANDLING =================

func (r *repository) DeleteBlocksFrom(height uint64) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Exec(
		`DELETE FROM blocks WHERE chain = ? AND height >= ?`,
		constants.ChainETH, height,
	).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Exec(
		`DELETE FROM transactions WHERE chain = ? AND block_height >= ?`,
		constants.ChainETH, height,
	).Error; err != nil {
		tx.Rollback()
		return err
	}

	newLast := uint64(0)
	if height > 0 {
		newLast = height - 1
	}

	if err := tx.Exec(
		`UPDATE indexer_states SET last_indexed_block = ? WHERE chain = ?`,
		newLast, constants.ChainETH,
	).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// ================= EXTENDED READ METHODS =================

func (r *repository) GetBlockByHash(hash string) (*model.Block, error) {
	var block model.Block
	err := r.db.Preload("Transactions").Where("chain = ? AND hash = ?", constants.ChainETH, hash).First(&block).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &block, nil
}

func (r *repository) GetTransactionByHash(hash string) (*model.Transaction, error) {
	var tx model.Transaction
	err := r.db.Where("chain = ? AND hash = ?", constants.ChainETH, hash).First(&tx).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tx, nil
}
func (r *repository) GetStats() (map[string]interface{}, error) {
	var blockCount int64
	var txCount int64

	if err := r.db.Model(&model.Block{}).Where("chain = ?", constants.ChainETH).Count(&blockCount).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&model.Transaction{}).Where("chain = ?", constants.ChainETH).Count(&txCount).Error; err != nil {
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
		constants.ChainETH, "%"+query+"%", query+"%").
		Order("height DESC").Limit(10).Find(&blocks).Error
	return blocks, err
}

func (r *repository) SearchTransactions(query string) ([]model.Transaction, error) {
	var txs []model.Transaction
	err := r.db.Where("chain = ? AND (hash ILIKE ? OR from_address ILIKE ? OR to_address ILIKE ?)",
		constants.ChainETH, "%"+query+"%", "%"+query+"%", "%"+query+"%").
		Order("block_height DESC").Limit(10).Find(&txs).Error
	return txs, err
}
