package handlers

import (
	"indexer/internal/model"
	"indexer/internal/repository"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type APIHandler struct {
	repo repository.Repository
}

func NewAPIHandler(repo repository.Repository) *APIHandler {
	return &APIHandler{repo: repo}
}

func (h *APIHandler) normalizeChain(c string) model.ChainType {
	switch strings.ToLower(c) {
	case "btc", "bitcoin":
		return model.ChainBTC
	case "eth", "ethereum":
		return model.ChainETH
	}
	return model.ChainType(c)
}

func (h *APIHandler) GetBlocks(c *gin.Context) {
	chain := h.normalizeChain(c.Param("chain"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	blocks, err := h.repo.GetLatestBlocks(chain, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch blocks"})
		return
	}

	total, _ := h.repo.CountBlocks(chain)

	dtos := make([]BlockResponse, len(blocks))
	for i, b := range blocks {
		dtos[i] = ToBlockDTO(b, int(b.TXCount))
	}

	c.JSON(http.StatusOK, PaginatedBlocksResponse{
		Page:   page,
		Limit:  limit,
		Total:  total,
		Blocks: dtos,
	})
}

func (h *APIHandler) GetBlockByHeight(c *gin.Context) {
	chain := h.normalizeChain(c.Param("chain"))
	height, err := strconv.ParseUint(c.Param("height"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid block height"})
		return
	}

	block, err := h.repo.GetBlockByHeight(chain, height)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Block not found"})
		return
	}

	txs, _ := h.repo.GetTransactionsByBlock(chain, height)

	txDTOs := make([]TransactionResponse, len(txs))
	for i, t := range txs {
		txDTOs[i] = ToTransactionDTO(t)
	}

	resp := BlockDetailsResponse{
		Height:       block.Height,
		Hash:         block.Hash,
		Timestamp:    block.Timestamp.Unix(),
		TxCount:      len(txs),
		Transactions: txDTOs,
	}

	c.JSON(http.StatusOK, resp)
}

func (h *APIHandler) GetStats(c *gin.Context) {
	btcLatest, _ := h.repo.GetState(model.ChainBTC)
	ethLatest, _ := h.repo.GetState(model.ChainETH)

	btcMax, _ := h.repo.GetMaxBlockHeight(model.ChainBTC)
	ethMax, _ := h.repo.GetMaxBlockHeight(model.ChainETH)

	btcBlocks, _ := h.repo.CountBlocks(model.ChainBTC)
	ethBlocks, _ := h.repo.CountBlocks(model.ChainETH)

	btcTxCount, _ := h.repo.CountTransactions(model.ChainBTC)
	ethTxCount, _ := h.repo.CountTransactions(model.ChainETH)

	// Sync logic: latestIndexedBlock >= (maxBlockInDB - 2)
	btcSynced := btcLatest >= (btcMax-2) && btcMax > 0
	ethSynced := ethLatest >= (ethMax-2) && ethMax > 0

	resp := StatsResponse{
		BTC: ChainStats{
			LatestBlock: btcLatest,
			TotalBlocks: btcBlocks,
			TotalTx:     btcTxCount,
			Synced:      btcSynced,
		},
		ETH: ChainStats{
			LatestBlock: ethLatest,
			TotalBlocks: ethBlocks,
			TotalTx:     ethTxCount,
			Synced:      ethSynced,
		},
	}

	c.JSON(http.StatusOK, resp)
}

func (h *APIHandler) Search(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Query parameter 'q' is required"})
		return
	}

	// 1. Try as block height (numeric)
	if height, err := strconv.ParseUint(q, 10, 64); err == nil {
		// Search BTC
		if block, err := h.repo.GetBlockByHeight(model.ChainBTC, height); err == nil {
			c.JSON(http.StatusOK, SearchResult{
				Type:   "block",
				Chain:  "btc",
				Result: ToBlockDTO(*block, 0),
			})
			return
		}
		// Search ETH
		if block, err := h.repo.GetBlockByHeight(model.ChainETH, height); err == nil {
			c.JSON(http.StatusOK, SearchResult{
				Type:   "block",
				Chain:  "eth",
				Result: ToBlockDTO(*block, 0),
			})
			return
		}
	}

	// 2. Try as transaction hash (hex string)
	if strings.HasPrefix(q, "0x") || len(q) >= 32 {
		// Search BTC
		if tx, err := h.repo.FindTransactionByHash(model.ChainBTC, q); err == nil {
			c.JSON(http.StatusOK, SearchResult{
				Type:   "transaction",
				Chain:  "btc",
				Result: ToTransactionDTO(*tx),
			})
			return
		}
		// Search ETH
		if tx, err := h.repo.FindTransactionByHash(model.ChainETH, q); err == nil {
			c.JSON(http.StatusOK, SearchResult{
				Type:   "transaction",
				Chain:  "eth",
				Result: ToTransactionDTO(*tx),
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, ErrorResponse{Error: "Not found"})
}
