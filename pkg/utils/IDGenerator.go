package utils

import (
    "crypto/rand"
    "fmt"
    "math/big"
    "log"
    "database/sql"
)

func GenerateRandomID(prefix string) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    var result = make([]byte, 16)

    for i := range result {
        randomByte, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
        if err != nil {
            log.Fatalf("Failed to generate random number: %v", err)
        }
        result[i] = charset[randomByte.Int64()]
    }

    return fmt.Sprintf(prefix + "%s", result)
}

func CheckUserIDIfExists(db *sql.DB, userID string) bool {
    var count int
    query := "SELECT COUNT(*) FROM users WHERE id = ?"
    err := db.QueryRow(query, userID).Scan(&count)
    if err != nil {
        log.Println("Error checking if ID exists:", err)
        return true
    }
    return count > 0
}

func CheckOrderIDIfExists(db *sql.DB, orderID string) bool {
    var count int
    query := "SELECT COUNT(*) FROM orders WHERE id = ?"
    err := db.QueryRow(query, orderID).Scan(&count)
    if err != nil {
        log.Println("Error checking if ID exists:", err)
        return true
    }
    return count > 0
}

func CheckItemIDIfExists(db *sql.DB, orderID string) bool {
    var count int
    query := "SELECT COUNT(*) FROM order_items WHERE order_item_id = ?"
    err := db.QueryRow(query, orderID).Scan(&count)
    if err != nil {
        log.Println("Error checking if ID exists:", err)
        return true
    }
    return count > 0
}