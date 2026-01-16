package eth

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	routes := router.Group("/eth")
	{
		routes.GET("/blocks", handler.GetBlocks)
		routes.GET("/transactions", handler.GetTransactions)
	}
}
