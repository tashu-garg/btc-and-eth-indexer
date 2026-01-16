package server

import (
	"indexer/internal/btc"
	"indexer/internal/eth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, ethHandler *eth.Handler, btcHandler *btc.Handler) {
	api := router.Group("/api/v1")
	{
		// ETH Routes
		ethRoutes := api.Group("/ethereum")
		{
			ethRoutes.GET("/blocks", ethHandler.GetBlocks)
			ethRoutes.GET("/blocks/:height", ethHandler.GetBlockByHeight)
			ethRoutes.GET("/block/hash/:hash", ethHandler.GetBlockByHash)
			ethRoutes.GET("/txs", ethHandler.GetTransactions)
			ethRoutes.GET("/tx/:hash", ethHandler.GetTransactionByHash)
			ethRoutes.GET("/stats", ethHandler.GetStats)
			ethRoutes.GET("/search", ethHandler.Search)
		}

		// BTC Routes
		btcRoutes := api.Group("/bitcoin")
		{
			btcRoutes.GET("/blocks", btcHandler.GetBlocks)
			btcRoutes.GET("/blocks/:height", btcHandler.GetBlockByHeight)
			btcRoutes.GET("/block/hash/:hash", btcHandler.GetBlockByHash)
			btcRoutes.GET("/txs", btcHandler.GetTransactions)
			btcRoutes.GET("/tx/:hash", btcHandler.GetTransactionByHash)
			btcRoutes.GET("/stats", btcHandler.GetStats)
			btcRoutes.GET("/search", btcHandler.Search)
		}
	}

	router.GET("/health", func(c *gin.Context) {
		// In a real app we might query services for their latest indexed block
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
}
