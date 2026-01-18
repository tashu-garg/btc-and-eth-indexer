package routes

import (
	"indexer/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(apiHandler *handlers.APIHandler) *gin.Engine {
	r := gin.Default()

	// CORS or other middleware can be added here
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := r.Group("/api")
	{
		api.GET("/stats", apiHandler.GetStats)
		api.GET("/search", apiHandler.Search)
		api.GET("/:chain/blocks", apiHandler.GetBlocks)
		api.GET("/:chain/blocks/:height", apiHandler.GetBlockByHeight)
	}

	return r
}
