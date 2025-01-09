package product_list

import (
	"database/sql"
	"slash-pos/internal/middleware"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(router *gin.Engine, db *sql.DB) {
	productRoutes := router.Group("/products")
	{
		productRoutes.Use(middleware.AuthMiddleware())
		productRoutes.GET("/", func(c *gin.Context) { GetProductsHandler(c, db) })         
		productRoutes.GET("/:id", func(c *gin.Context) { SearchProductByIDhandler(c, db) })        
	}
}
