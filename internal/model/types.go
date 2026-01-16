package model

import (
	"time"
)

type ChainType string

const (
	ChainBTC ChainType = "bitcoin"
	ChainETH ChainType = "ethereum"
)

// Block represents a generic blockchain block
type Block struct {
	ID           uint64        `json:"id" gorm:"primaryKey;autoIncrement"`
	Chain        ChainType     `json:"chain" gorm:"type:varchar(10);uniqueIndex:idx_chain_height"`
	Height       uint64        `json:"height" gorm:"uniqueIndex:idx_chain_height"`
	Hash         string        `json:"hash" gorm:"uniqueIndex"`
	BlockHash    string        `json:"block_hash" gorm:"index"`
	Transactions []Transaction `json:"transactions,omitempty" gorm:"-"`
	Timestamp    time.Time     `json:"timestamp"`
	CreatedAt    time.Time     `json:"created_at"`
}

// Transaction represents a generic blockchain transaction
type Transaction struct {
	ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	Chain     ChainType `json:"chain" gorm:"type:varchar(10);index"`
	Hash      string    `json:"hash" gorm:"uniqueIndex"`
	BlockHash string    `json:"block_hash" gorm:"index"`
	Height    uint64    `json:"height" gorm:"column:block_height;index"` // Match existing column name
	From      string    `json:"from_address" gorm:"column:from_address"`
	To        string    `json:"to_address" gorm:"column:to_address"`
	Value     string    `json:"value"`  // string to handle big numbers (Satoshis/Wei)
	Status    string    `json:"status"` // "success", "failed", "pending"
	Timestamp time.Time `json:"timestamp"`
	CreatedAt time.Time `json:"created_at"`
}

// IndexerState tracks the indexing progress
type IndexerState struct {
	Chain             ChainType `json:"chain" gorm:"primaryKey;type:varchar(10)"`
	LastIndexedHeight uint64    `json:"last_indexed_height" gorm:"column:last_indexed_block"`
	UpdatedAt         time.Time `json:"updated_at"`
}
