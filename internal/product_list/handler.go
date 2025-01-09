package product_list

import (
	"slash-pos/pkg/logger"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var log = logger.NewLogger(true)

// Struct untuk produk
type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
	Size  string `json:"size"`
	Stock int    `json:"stock"`
	Prize string `json:"prize"`
}

func GetProductsHandler(c *gin.Context, db *sql.DB) {
	query := "SELECT id, name, color, size, stock, price FROM products"

	rows, err := db.Query(query)
	if err != nil {
		log.Error("Failed to fetch products:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Failed to fetch products",
			"error":   err.Error(),
		})
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Color, &product.Size, &product.Stock, &product.Prize); err != nil {
			log.Error("Failed to scan product:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "failed",
				"message": "Failed to scan product",
				"error":   err.Error(),
			})
			return
		}
		products = append(products, product)
	}

	log.Info("Successfully fetched products")
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Products fetched successfully",
		"data":    products,
	})
}


func SearchProductByIDhandler(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	query := "SELECT id, name, color, size, stock, price FROM products WHERE id = ?"

	var product Product
	err := db.QueryRow(query, id).Scan(&product.ID, &product.Name, &product.Color, &product.Size, &product.Stock, &product.Prize)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": "Product not found",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Failed to fetch product",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Product fetched successfully",
		"data":    product,
	})
}