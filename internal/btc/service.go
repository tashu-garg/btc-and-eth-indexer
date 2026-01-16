package btc

import (
	"context"
	"encoding/json"
	"fmt"
	"indexer/internal/config"
	"indexer/internal/constants"
	"indexer/internal/model"
	"log"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
)

type BTCClient struct {
	Client  *http.Client
	url     string
	user    string
	pass    string
	timeout time.Duration
}

type Service struct {
	client *BTCClient
	repo   Repository
	cfg    config.BTCConfig
}

func NewService(db *gorm.DB, cfg config.BTCConfig) (*Service, error) {
	return &Service{
		client: NewBTCClient(&cfg),
		repo:   NewRepository(db),
		cfg:    cfg,
	}, nil
}

func NewBTCClient(cfg *config.BTCConfig) *BTCClient {
	if cfg.RPCURL == "" {
		log.Println("WARNING: BTC_RPC_URL is not set")
	}
	return &BTCClient{
		Client:  &http.Client{Timeout: 30 * time.Second},
		url:     cfg.RPCURL,
		user:    cfg.RPCUser,
		pass:    cfg.RPCPass,
		timeout: 30 * time.Second,
	}
}

func (c *BTCClient) URL() string {
	return c.url
}

// BTC Block structure from RPC
type BlockbookStatus struct {
	Backend struct {
		Blocks uint64 `json:"blocks"`
	} `json:"backend"`
}

type RPCBlock struct {
	Hash              string           `json:"hash"`
	Height            uint64           `json:"height"`
	PreviousBlockHash string           `json:"previousBlockHash"`
	Time              int64            `json:"time"`
	Txs               []RPCTransaction `json:"txs"`
}

type RPCTransaction struct {
	TxID string `json:"txid"`
	Vin  []struct {
		Addresses []string `json:"addresses"`
		Value     string   `json:"value"`
	} `json:"vin"`
	Vout []struct {
		Value     string   `json:"value"`
		Addresses []string `json:"addresses"`
	} `json:"vout"`
	Value string `json:"value"`
}

func (s *Service) Start(ctx context.Context) {
	log.Println("BTC sync loop started")
	// No continuous loop here if we want to rely on the cron scheduler,
	// but the user's Go code had a Run/Start loop.
	// I'll keep it as a wait-loop like ETH to avoid double-processing.
	<-ctx.Done()
}

func (s *Service) ProcessNextBlock(ctx context.Context) error {
	lastIndexed, err := s.repo.GetLastIndexedBlock()
	if err != nil {
		return err
	}

	// Fast-forward if starting from scratch to show data immediately
	if lastIndexed == 0 {
		lastIndexed = 118820
	}

	nextHeight := lastIndexed + 1
	if s.client.url == "" {
		return nil
	}

	tip, err := s.client.GetBlockCount(ctx)
	if err != nil {
		log.Printf("[BTC_DEBUG] Failed to get tip: %v", err)
		return err
	}

	if nextHeight > tip {
		return nil
	}

	log.Printf("[BTC_DEBUG] Indexing block %d / %d", nextHeight, tip)
	block, txs, err := s.client.GetBlock(ctx, nextHeight)
	if err != nil {
		return err
	}

	if err := s.repo.SaveBlock(block); err != nil {
		return err
	}

	for _, tx := range txs {
		if err := s.repo.SaveTransaction(tx); err != nil {
			return err
		}
	}

	if err := s.repo.UpdateIndexerState(nextHeight); err != nil {
		return err
	}

	log.Printf("Indexed BTC block %d (%d txs)", nextHeight, len(txs))
	return nil
}

func (c *BTCClient) GetBlockCount(ctx context.Context) (uint64, error) {
	url := fmt.Sprintf("%s/api/v2/status?apikey=%s", strings.TrimSuffix(c.url, "/"), c.pass)
	if !strings.Contains(c.url, "?") && c.pass == "" {
		// Fallback if apikey is already in URL
		url = fmt.Sprintf("%s/api/v2/status", strings.TrimSuffix(c.url, "/"))
	} else if strings.Contains(c.url, "apikey=") {
		url = strings.Replace(c.url, "/?", "/api/v2/status?", 1)
		if !strings.Contains(url, "/api/v2/status") {
			url = fmt.Sprintf("%s/api/v2/status", strings.TrimSuffix(c.url, "/"))
		}
	}

	// Simplest approach: manual URL construction based on seen pattern
	baseUrl := strings.Split(c.url, "?")[0]
	apiKey := ""
	if parts := strings.Split(c.url, "apikey="); len(parts) > 1 {
		apiKey = parts[1]
	}
	url = fmt.Sprintf("%s/api/v2/status?apikey=%s", strings.TrimSuffix(baseUrl, "/"), apiKey)

	resp, err := c.Client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var status BlockbookStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return 0, err
	}
	return status.Backend.Blocks, nil
}

func (c *BTCClient) GetBlock(ctx context.Context, height uint64) (*model.Block, []*model.Transaction, error) {
	baseUrl := strings.Split(c.url, "?")[0]
	apiKey := ""
	if parts := strings.Split(c.url, "apikey="); len(parts) > 1 {
		apiKey = parts[1]
	}
	url := fmt.Sprintf("%s/api/v2/block/%d?apikey=%s", strings.TrimSuffix(baseUrl, "/"), height, apiKey)

	resp, err := c.Client.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	var rpcBlock RPCBlock
	if err := json.NewDecoder(resp.Body).Decode(&rpcBlock); err != nil {
		return nil, nil, err
	}

	block := &model.Block{
		Chain:     constants.ChainBTC,
		Height:    rpcBlock.Height,
		Hash:      rpcBlock.Hash,
		BlockHash: rpcBlock.PreviousBlockHash,
		Timestamp: time.Unix(rpcBlock.Time, 0),
	}

	var txs []*model.Transaction
	for _, rawTx := range rpcBlock.Txs {
		fromAddress := ""
		if len(rawTx.Vin) > 0 && len(rawTx.Vin[0].Addresses) > 0 {
			fromAddress = rawTx.Vin[0].Addresses[0]
		}

		toAddress := ""
		if len(rawTx.Vout) > 0 && len(rawTx.Vout[0].Addresses) > 0 {
			toAddress = rawTx.Vout[0].Addresses[0]
		}

		// Blockbook value is in Satoshis as string
		valSat := 0.0
		fmt.Sscanf(rawTx.Value, "%f", &valSat)
		valBTC := valSat / 100000000.0

		tx := &model.Transaction{
			Chain:     constants.ChainBTC,
			Hash:      rawTx.TxID,
			BlockHash: rpcBlock.Hash,
			Height:    rpcBlock.Height,
			From:      fromAddress,
			To:        toAddress,
			Value:     fmt.Sprintf("%.8f", valBTC),
			Status:    "success",
			Timestamp: block.Timestamp,
		}
		txs = append(txs, tx)
	}

	return block, txs, nil
}

// API METHODS
func (s *Service) GetBlocks(limit, offset int) ([]model.Block, error) {
	return s.repo.GetBlocks(limit, offset)
}

func (s *Service) GetTransactions(limit, offset int) ([]model.Transaction, error) {
	return s.repo.GetTransactions(limit, offset)
}

func (s *Service) GetBlockByHeight(height uint64) (*model.Block, error) {
	return s.repo.GetBlockByHeight(height)
}

func (s *Service) GetBlockByHash(hash string) (*model.Block, error) {
	return s.repo.GetBlockByHash(hash)
}

func (s *Service) GetTransactionByHash(hash string) (*model.Transaction, error) {
	return s.repo.GetTransactionByHash(hash)
}

func (s *Service) GetStats() (map[string]interface{}, error) {
	return s.repo.GetStats()
}

func (s *Service) SearchBlocks(query string) ([]model.Block, error) {
	return s.repo.SearchBlocks(query)
}

func (s *Service) SearchTransactions(query string) ([]model.Transaction, error) {
	return s.repo.SearchTransactions(query)
}
