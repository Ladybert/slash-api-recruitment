package main

import (
	"slash-pos/config"
	"slash-pos/pkg/db"
	"slash-pos/pkg/logger"
	"slash-pos/internal/authentication"
	"slash-pos/internal/product_list"
	"slash-pos/internal/order"

	"github.com/gin-gonic/gin"
)

func main() {
	appConfig := config.LoadConfig()

	log := logger.NewLogger(true)
	log.Info("Starting application...")

	database := db.InitDB(appConfig.DB)
	defer database.Close()
	log.Info("Connected to the database")

	router := gin.Default()
	authentication.RegisterRoutes(router, database)
	product_list.ProductRoutes(router, database)
	order.OrderRoutes(router, database)

	// Run server
	log.Info("Starting server on", appConfig.Server.Address)
	if err := router.Run(appConfig.DB.Host + appConfig.Server.Address); err != nil {
		log.Error("Failed to start server:", err)
	}
}
