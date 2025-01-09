	package order

	import (
		"database/sql"
		"fmt"
		"net/http"
		"slash-pos/pkg/logger"
		"slash-pos/pkg/utils"
		"strconv"
		"time"
		"strings"

		"github.com/gin-gonic/gin"
	)

	var log = logger.NewLogger(true)

	func formatToRupiah(amount int) string {
		if amount < 0 {
			return "Rp 0"
		}
	
		strAmount := strconv.Itoa(amount)
	
		var result []string
		for i := len(strAmount); i > 0; i -= 3 {
			start := i - 3
			if start < 0 {
				start = 0
			}
			result = append([]string{strAmount[start:i]}, result...)
		}
	
		return "Rp " + strings.Join(result, ".")
	}
	

	func createOrderHandler(c *gin.Context, db *sql.DB) {
		var newOrder Order
		var prefix = "order-"
	
		if err := c.ShouldBindJSON(&newOrder); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newOrderTime := time.Now()
		newOrder.Status = "pending"
		expiredTime := newOrderTime.Add(1 * time.Hour)
		formattedExpiredTime := expiredTime.Format("2006-01-02 15:04:05")
	
		orderID := utils.GenerateRandomID(prefix)
		for utils.CheckOrderIDIfExists(db, orderID) {
			orderID = utils.GenerateRandomID(prefix)
		}
		newOrder.ID = orderID
	
		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
			return
		}
	
		totalCalculated := 0

		formattedOrderTime := newOrderTime.Format("2006-01-02 15:04:05")
		
	
		log.Info(newOrder.ID, newOrder.UserID, newOrder.CustomerName, newOrder.CustomerPhone, newOrder.CustomerAddress, "Rp 0", newOrder.OrderTime, newOrder.Expired, newOrder.Status)

		queryOrder := "INSERT INTO orders (id, user_id, customer_name, customer_phone, customer_address, total_amount, order_time, expired, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
		_, err = tx.Exec(queryOrder, newOrder.ID, newOrder.UserID, newOrder.CustomerName, newOrder.CustomerPhone, newOrder.CustomerAddress, "Rp 0", formattedOrderTime, formattedExpiredTime, newOrder.Status)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
			return
		}
	
		for _, item := range newOrder.OrderItems {
			if item.Quantity <= 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order input"})
				tx.Rollback()
				return
			}
	
			var stock int
			var price string
			queryStock := "SELECT stock, price FROM products WHERE id = ?"
			err := tx.QueryRow(queryStock, item.ProductID).Scan(&stock, &price)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product stock"})
				return
			}
	
			if item.Quantity > stock {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf("Stock produk tidak mencukupi jumlah pesanan produk: %s", item.ProductID),
				})
				tx.Rollback()
				return
			}
	
			newStock := stock - item.Quantity
			queryUpdateStock := "UPDATE products SET stock = ? WHERE id = ?"
			_, err = tx.Exec(queryUpdateStock, newStock, item.ProductID)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product stock"})
				return
			}
	
			priceCleaned := strings.ReplaceAll(price, "Rp ", "")
			priceCleaned = strings.ReplaceAll(priceCleaned, ".", "")
			priceInt, err := strconv.Atoi(priceCleaned)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid price format in database"})
				return
			}
	
			itemTotal := priceInt * item.Quantity
			totalCalculated += itemTotal
	
			item.Price = formatToRupiah(priceInt)
			itemTotalFormatted := formatToRupiah(itemTotal)
	
			fmt.Println("Item Total (Formatted):", itemTotalFormatted)
	
			itemID := utils.GenerateRandomID("oid-")
			for utils.CheckItemIDIfExists(db, itemID) {
				itemID = utils.GenerateRandomID("oid-")
			}
	
			queryItem := "INSERT INTO order_items (order_item_id, order_id, product_id, quantity, price) VALUES (?, ?, ?, ?, ?)"
			_, err = tx.Exec(queryItem, itemID, newOrder.ID, item.ProductID, item.Quantity, item.Price)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert order items"})
				return
			}
		}

		queryUpdateOrder := "UPDATE orders SET total_amount = ? WHERE id = ?"
		formattedAmount := formatToRupiah(totalCalculated)
		_, err = tx.Exec(queryUpdateOrder, formattedAmount, newOrder.ID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update total amount in order"})
			return
		}
	
		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}
	
		c.JSON(http.StatusCreated, gin.H{
			"expiredTime": formattedExpiredTime,
			"orderID":     newOrder.ID,
			"message":     "Order created successfully",
			"status":      "success",
		})
	}

	func getDetailOrderByIDHandler(c *gin.Context, db *sql.DB) {
		orderID := c.Param("id")
	
		getDetailQuery := `
		SELECT o.id, o.customer_name, o.customer_phone, o.customer_address, o.total_amount, o.order_time, o.expired, o.status,
				oi.product_id, oi.quantity, oi.price
		FROM orders o
		LEFT JOIN order_items oi ON o.id = oi.order_id
		WHERE o.id = ?`
	
		rows, err := db.Query(getDetailQuery, orderID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve order details"})
			return
		}
		defer rows.Close()
	
		var isOrderProcessed bool
		var orderDetails OrderDetails
	
		for rows.Next() {
			var item OrderItem
	
			err := rows.Scan(
				&orderDetails.ID,
				&orderDetails.CustomerName,
				&orderDetails.CustomerPhone,
				&orderDetails.CustomerAddress,
				&orderDetails.TotalAmount,
				&orderDetails.OrderTime,
				&orderDetails.Expired,
				&orderDetails.Status,
				&item.ProductID,
				&item.Quantity,
				&item.Price,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
				return
			}
	
			if !isOrderProcessed {
				isOrderProcessed = true
			}
	
			orderDetails.Items = append(orderDetails.Items, item)

		}
	
		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while scanning rows"})
			return
		}
	
		c.JSON(http.StatusOK, orderDetails)
	}

	func payNowOrderByIDHandler(c *gin.Context, db *sql.DB) {
		var order OrderDetails
	
		var orderIDForm struct {
			OrderID string `json:"orderID"`
		}
	
		if err := c.ShouldBindJSON(&orderIDForm); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	
		getOrderQuery := "SELECT expired, status FROM orders WHERE id = ?"
		err := db.QueryRow(getOrderQuery, orderIDForm.OrderID).Scan(&order.Expired, &order.Status)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve order"})
			}
			return
		}

		loc, err := time.LoadLocation("Asia/Jakarta") // WIB
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load time location"})
			return
		}
	
		expiredTimeParsed, err := time.Parse("2006-01-02 15:04:05", order.Expired)
		if err != nil {
			expiredTimeParsed, err = time.ParseInLocation("2006-01-02 15:04:05", order.Expired, time.UTC)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse expired time"})
				return
			}
		}

		currentTime := time.Now().In(loc)

		log.Info(currentTime)

		if currentTime.After(expiredTimeParsed) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Waktu pesanan sudah kadaluwarsa"})
			return
		}
	
		if order.Status == "done" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Pesanan sudah dibayar"})
			return
		}
	
		updateStatusQuery := "UPDATE orders SET status = 'done' WHERE id = ?"
		_, err = db.Exec(updateStatusQuery, orderIDForm.OrderID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
			return
		}
	
		c.JSON(http.StatusOK, gin.H{"message": "Payment successful, order status updated to done"})
	}
	
	func updateOrderByIDHandler(c *gin.Context, db *sql.DB) {
		orderID := c.Param("id")
	
		var newOrder struct {
			ProductID string `json:"productID"`
			Quantity  int    `json:"quantity"`
		}
	
		if err := c.ShouldBindJSON(&newOrder); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	
		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
			return
		}
		defer tx.Rollback()
	
		var currentItem struct {
			Quantity int    `json:"quantity"`
			ProductID string `json:"product_id"`
			Price     string    `json:"price"`
		}
		queryGetOrderItem := "SELECT quantity, product_id, price FROM order_items WHERE order_id = ?"
		err = tx.QueryRow(queryGetOrderItem, orderID).Scan(&currentItem.Quantity, &currentItem.ProductID, &currentItem.Price)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Order item not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve order item"})
			}
			return
		}
	
		var productPrice string
		queryGetProductPrice := "SELECT price, stock FROM products WHERE id = ?"
		var currentStock int
		err = tx.QueryRow(queryGetProductPrice, newOrder.ProductID).Scan(&productPrice, &currentStock)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product price"})
			return
		}
	
		if newOrder.Quantity > currentItem.Quantity {
			difference := newOrder.Quantity - currentItem.Quantity
			if currentStock < difference {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
				return
			}

			_, err = tx.Exec("UPDATE products SET stock = stock - ? WHERE id = ?", difference, newOrder.ProductID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product stock"})
				return
			}
		} else if newOrder.Quantity < currentItem.Quantity {
			difference := currentItem.Quantity - newOrder.Quantity
			_, err = tx.Exec("UPDATE products SET stock = stock + ? WHERE id = ?", difference, newOrder.ProductID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product stock"})
				return
			}
		}

		cleanedPrice := strings.ReplaceAll(productPrice, "Rp ", "")
		cleanedPrice = strings.ReplaceAll(cleanedPrice, ".", "") 

		productPriceInt, err := strconv.Atoi(cleanedPrice)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert product price to integer"})
			return
		}

		priceUpdate := newOrder.Quantity * productPriceInt
		formattedPriceUpdate := formatToRupiah(priceUpdate)

		queryUpdateOrderItem := `
			UPDATE order_items 
			SET quantity = ?, price = ? 
			WHERE order_id = ? AND product_id = ?`
		_, err = tx.Exec(queryUpdateOrderItem, newOrder.Quantity, formattedPriceUpdate, orderID, newOrder.ProductID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order item"})
			return
		}
	
		var updatedOrderAmount int
		queryGetTotalAmount := `
			SELECT oi.price, oi.quantity 
			FROM order_items oi
			WHERE oi.order_id = ?`
		rows, err := tx.Query(queryGetTotalAmount, orderID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve order items"})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var price string
			var quantity int
			if err := rows.Scan(&price, &quantity); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan order item"})
				return
			}

			priceInt, err := strconv.Atoi(strings.ReplaceAll(strings.ReplaceAll(price, "Rp ", ""), ".", ""))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert price"})
				return
			}

			log.Info(priceInt)

			updatedOrderAmount += priceInt 
		}

		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over order items"})
			return
		}

		log.Info(updatedOrderAmount)

		formattedAmount := formatToRupiah(updatedOrderAmount)
	
		queryUpdateTotalAmount := "UPDATE orders SET total_amount = ? WHERE id = ?"
		_, err = tx.Exec(queryUpdateTotalAmount, formattedAmount, orderID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update total amount"})
			return
		}
	
		err = tx.Commit()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}
	
		c.JSON(http.StatusOK, gin.H{"message": "Order updated successfully"})
	}

	func deleteOrderByIDHandler(c *gin.Context, db *sql.DB) {
		orderID := c.Param("id")
	
		userID := c.GetString("userID")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
	
		var orderOwnerID string
		queryCheckOrder := "SELECT user_id FROM orders WHERE id = ?"
		err := db.QueryRow(queryCheckOrder, orderID).Scan(&orderOwnerID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve order"})
			}
			return
		}
	
		if orderOwnerID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this order"})
			return
		}
	
		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
			return
		}
		defer tx.Rollback()
	
		_, err = tx.Exec("DELETE FROM order_items WHERE order_id = ?", orderID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order items"})
			return
		}
	
		_, err = tx.Exec("DELETE FROM orders WHERE id = ?", orderID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order"})
			return
		}
	
		err = tx.Commit()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}
	
		c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
	}
	