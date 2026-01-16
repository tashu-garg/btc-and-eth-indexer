package eth

import (
	"context"
	"fmt"
	"indexer/internal/config"
	"indexer/internal/constants"
	"indexer/internal/model"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Service struct {
	client *ethclient.Client
	repo   Repository

	mu sync.Mutex

	chainID *big.Int // ✅ add this

	lastKnownTip uint64
	lastTipCheck time.Time

	logger *logrus.Entry
	cfg    config.ETHConfig
}

func NewService(db *gorm.DB, cfg config.ETHConfig) (*Service, error) {
	client, err := ethclient.Dial(cfg.RPCURL)
	if err != nil {
		return nil, err
	}

	svc := &Service{
		client: client,
		repo:   NewRepository(db),
		logger: logrus.WithField("chain", constants.ChainETH),
		cfg:    cfg,
	}

	// Try to get chainID, but don't hard fail if node is temporarily flaky
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chainID, err := client.NetworkID(ctx)
	if err == nil {
		svc.chainID = chainID
	} else {
		svc.logger.Warnf("Could not fetch ChainID during initialization (node may be flaky): %v. Will retry during processing.", err)
	}

	return svc, nil
}

// Start must be BLOCKING
// Start is now just a logger, as the Cron scheduler drives the execution
// Run starts the service but now relies on the scheduler in app.go to call ProcessNextBlock
func (s *Service) Run(ctx context.Context) {
	s.logger.Info("ETH sync service started (waiting for scheduler)")
	<-ctx.Done()
	s.logger.Info("ETH sync service stopping")
}
func isRateLimit(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "429") ||
		strings.Contains(msg, "Too Many Requests") ||
		strings.Contains(msg, "rate limit")
}

// ProcessNextBlock processes a single block. It is called by the external scheduler.
func (s *Service) ProcessNextBlock() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx := context.Background()

	lastIndexed, err := s.repo.GetLastIndexedBlock()
	if err != nil {
		return err
	}

	nextHeight := lastIndexed + 1

	// refresh chain tip every 5s only
	if time.Since(s.lastTipCheck) > 5*time.Second {
		tip, err := s.client.BlockNumber(ctx)
		if err != nil {
			return err
		}
		s.lastKnownTip = tip
		s.lastTipCheck = time.Now()
	}

	if nextHeight > s.lastKnownTip {
		time.Sleep(3 * time.Second)
		return nil
	}

	block, err := s.client.BlockByNumber(ctx, big.NewInt(int64(nextHeight)))
	if err != nil {
		return err
	}

	modelBlock := &model.Block{
		Chain:     constants.ChainETH,
		Height:    block.NumberU64(),
		Hash:      block.Hash().Hex(),
		BlockHash: block.ParentHash().Hex(),
		Timestamp: time.Unix(int64(block.Time()), 0),
	}

	if err := s.repo.SaveBlock(modelBlock); err != nil {
		return err
	}

	var txs []model.Transaction

	// Ensure chainID is loaded
	if s.chainID == nil {
		chainID, err := s.client.NetworkID(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch chainID: %w", err)
		}
		s.chainID = chainID
	}

	signer := types.LatestSignerForChainID(s.chainID)

	for _, tx := range block.Transactions() {
		from := ""
		if sender, err := types.Sender(signer, tx); err == nil {
			from = sender.Hex()
		}

		to := ""
		if tx.To() != nil {
			to = tx.To().Hex()
		}

		txs = append(txs, model.Transaction{
			Hash:   tx.Hash().Hex(),
			Chain:  constants.ChainETH,
			Height: block.NumberU64(),
			From:   from,
			To:     to,
			Value:  tx.Value().String(),
		})
	}

	if len(txs) > 0 {
		if err := s.repo.SaveTransactions(txs); err != nil {
			return err
		}
	}

	if err := s.repo.UpdateIndexerState(block.NumberU64()); err != nil {
		return err
	}

	s.logger.Infof("Indexed ETH block %d (%d txs)", block.NumberU64(), len(txs))
	return nil
}

// ================= API METHODS =================

func (s *Service) GetBlocks(limit, offset int) ([]model.Block, error) {
	return s.repo.GetBlocks(limit, offset)
}

func (s *Service) GetTransactions(limit, offset int) ([]model.Transaction, error) {
	return s.repo.GetTransactions(limit, offset)
}

// ================= EXTENDED READ METHODS =================

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
