package workers

import (
	"context"
	"fmt"
	"indexer/internal/model"
	"indexer/internal/repository"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ETHWorker struct {
	repo         repository.Repository
	client       *ethclient.Client
	startHeight  int
	syncInterval time.Duration
}

func NewETHWorker(repo repository.Repository, rpcURL string, startHeight, syncIntervalMS int) (*ETHWorker, error) {
	if rpcURL == "" {
		return nil, fmt.Errorf("ETH_RPC_URL not configured")
	}
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	return &ETHWorker{
		repo:         repo,
		client:       client,
		startHeight:  startHeight,
		syncInterval: time.Duration(syncIntervalMS) * time.Millisecond,
	}, nil
}

func (w *ETHWorker) Start(ctx context.Context) {
	log.Println("[ETH] Worker starting...")

	// 1. Validate RPC connection
	tip, err := w.client.BlockNumber(ctx)
	if err != nil {
		log.Fatalf("[ETH] RPC connection failed: %v", err)
	}
	log.Printf("[ETH] RPC connection validated. Current tip: %d", tip)

	// 2. Ensure state exists
	_, err = w.repo.GetOrCreateState(model.ChainETH, tip, w.startHeight)
	if err != nil {
		log.Fatalf("[ETH] Failed to initialize state: %v", err)
	}

	ticker := time.NewTicker(w.syncInterval)
	defer ticker.Stop()

	log.Println("[ETH] Worker sync loop started")
	for {
		select {
		case <-ctx.Done():
			log.Println("[ETH] Worker stopping...")
			return
		case <-ticker.C:
			if err := w.sync(ctx); err != nil {
				log.Printf("[ETH] Sync error: %v", err)
			}
		}
	}
}

func (w *ETHWorker) sync(ctx context.Context) error {
	tip, err := w.client.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get tip: %w", err)
	}

	lastIndexed, err := w.repo.GetOrCreateState(model.ChainETH, tip, w.startHeight)
	if err != nil {
		return fmt.Errorf("failed to get state: %w", err)
	}

	if lastIndexed >= tip {
		return nil
	}

	nextHeight := lastIndexed + 1

	log.Printf("[ETH] Syncing block %d / %d", nextHeight, tip)

	block, err := w.client.BlockByNumber(ctx, big.NewInt(int64(nextHeight)))
	if err != nil {
		return fmt.Errorf("failed to fetch block %d: %w", nextHeight, err)
	}

	modelBlock := &model.Block{
		Chain:     model.ChainETH,
		Height:    block.NumberU64(),
		Hash:      block.Hash().Hex(),
		BlockHash: block.ParentHash().Hex(),
		TXCount:   uint64(len(block.Transactions())),
		Timestamp: time.Unix(int64(block.Time()), 0),
	}

	chainID, _ := w.client.NetworkID(ctx)
	signer := types.LatestSignerForChainID(chainID)

	var txs []*model.Transaction
	for _, tx := range block.Transactions() {
		from := ""
		if sender, err := types.Sender(signer, tx); err == nil {
			from = sender.Hex()
		}
		to := ""
		if tx.To() != nil {
			to = tx.To().Hex()
		}

		txs = append(txs, &model.Transaction{
			Chain:     model.ChainETH,
			Hash:      tx.Hash().Hex(),
			BlockHash: block.Hash().Hex(),
			Height:    block.NumberU64(),
			From:      from,
			To:        to,
			Value:     tx.Value().String(),
			Status:    "success",
			Timestamp: modelBlock.Timestamp,
		})
	}

	return w.repo.SaveBlockWithTransactions(modelBlock, txs)
}
