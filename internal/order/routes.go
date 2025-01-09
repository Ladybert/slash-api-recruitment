package order

import (
	"database/sql"
	"slash-pos/internal/middleware"

	"github.com/gin-gonic/gin"
)

func OrderRoutes(router *gin.Engine, db *sql.DB) {
	order := router.Group("/order")
	{
		order.Use(middleware.AuthMiddleware())
		order.POST("/", func(c *gin.Context) {
			createOrderHandler(c, db)
		})
		order.GET("/:id", func(c *gin.Context) {
			getDetailOrderByIDHandler(c, db)
		})
		order.POST("/pay-now", func(c *gin.Context) {
			payNowOrderByIDHandler(c, db)
		})
		order.PUT("/update/:id", func(c *gin.Context) {
			updateOrderByIDHandler(c, db)
		})
		order.DELETE("/order-delete/:id", func(c *gin.Context) {
			deleteOrderByIDHandler(c, db)
		})
	}
}
