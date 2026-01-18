package handlers

import "indexer/internal/model"

type BlockResponse struct {
	Height    uint64 `json:"height"`
	Hash      string `json:"hash"`
	TxCount   int    `json:"txCount"`
	Timestamp int64  `json:"timestamp"`
}

type TransactionResponse struct {
	Hash      string `json:"hash"`
	From      string `json:"from"`
	To        string `json:"to"`
	Value     string `json:"value"`
	Height    uint64 `json:"height"`
	Timestamp int64  `json:"timestamp"`
}

type PaginatedBlocksResponse struct {
	Page   int             `json:"page"`
	Limit  int             `json:"limit"`
	Total  int64           `json:"total"`
	Blocks []BlockResponse `json:"blocks"`
}

type BlockDetailsResponse struct {
	Height       uint64                `json:"height"`
	Hash         string                `json:"hash"`
	Timestamp    int64                 `json:"timestamp"`
	TxCount      int                   `json:"txCount"`
	Transactions []TransactionResponse `json:"transactions"`
}

type SearchResult struct {
	Type   string      `json:"type"` // "block" or "transaction"
	Chain  string      `json:"chain"`
	Result interface{} `json:"result"`
}

type StatsResponse struct {
	BTC ChainStats `json:"btc"`
	ETH ChainStats `json:"eth"`
}

type ChainStats struct {
	LatestBlock uint64 `json:"latestBlock"`
	TotalBlocks int64  `json:"totalBlocks"`
	TotalTx     int64  `json:"totalTx"`
	Synced      bool   `json:"synced"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func ToBlockDTO(b model.Block, txCount int) BlockResponse {
	return BlockResponse{
		Height:    b.Height,
		Hash:      b.Hash,
		TxCount:   txCount,
		Timestamp: b.Timestamp.Unix(),
	}
}

func ToTransactionDTO(t model.Transaction) TransactionResponse {
	return TransactionResponse{
		Hash:      t.Hash,
		From:      t.From,
		To:        t.To,
		Value:     t.Value,
		Height:    t.Height,
		Timestamp: t.Timestamp.Unix(),
	}
}
