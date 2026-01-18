package model

import (
	"time"
)

type ChainType string

const (
	ChainBTC ChainType = "bitcoin"
	ChainETH ChainType = "ethereum"
)

// Block represents the shared block structure
type Block struct {
	ID           uint64        `json:"id" gorm:"primaryKey;autoIncrement"`
	Chain        ChainType     `json:"chain" gorm:"type:varchar(10);uniqueIndex:idx_chain_height"`
	Height       uint64        `json:"height" gorm:"uniqueIndex:idx_chain_height"`
	Hash         string        `json:"hash" gorm:"index"`
	BlockHash    string        `json:"block_hash" gorm:"index"` // This is usually parent hash
	Transactions []Transaction `json:"transactions,omitempty" gorm:"-"`
	TXCount      uint64        `json:"tx_count"`
	Timestamp    time.Time     `json:"timestamp"`
	CreatedAt    time.Time     `json:"created_at"`
}

// Transaction represents the shared transaction structure
type Transaction struct {
	ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	Chain     ChainType `json:"chain" gorm:"type:varchar(10);index"`
	Hash      string    `json:"hash" gorm:"index"`
	BlockHash string    `json:"block_hash" gorm:"index"`
	Height    uint64    `json:"height" gorm:"column:block_height;index"`
	From      string    `json:"from_address" gorm:"column:from_address"`
	To        string    `json:"to_address" gorm:"column:to_address"`
	Value     string    `json:"value"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	CreatedAt time.Time `json:"created_at"`
}

// IndexerState tracks the indexing progress
type IndexerState struct {
	Chain             ChainType `json:"chain" gorm:"primaryKey;type:varchar(10)"`
	LastIndexedHeight uint64    `json:"last_indexed_height" gorm:"column:last_indexed_block"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// Table definitions for GORM migration
type BTCBlock struct{ Block }

func (BTCBlock) TableName() string { return "btc_blocks" }

type BTCTransaction struct{ Transaction }

func (BTCTransaction) TableName() string { return "btc_transactions" }

type ETHBlock struct{ Block }

func (ETHBlock) TableName() string { return "eth_blocks" }

type ETHTransaction struct{ Transaction }

func (ETHTransaction) TableName() string { return "eth_transactions" }
