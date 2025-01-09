package authentication

import (
	"slash-pos/pkg/logger"
	"database/sql"
	"net/http"
	"slash-pos/pkg/utils"

	"github.com/gin-gonic/gin"
)

var log = logger.NewLogger(true)
    
func RegisterUserHandler(c *gin.Context, db *sql.DB) {
    var newUser User
	var prefix = "user-"

    if err := c.ShouldBindJSON(&newUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
			"status": "Failed",
			"error": "Invalid input",
		})
        return
    }

    userID := utils.GenerateRandomID(prefix)

    for utils.CheckUserIDIfExists(db, userID) {
        userID = utils.GenerateRandomID(prefix)
    }

    newUser.ID = userID

    hashedPassword, err := utils.HashPassword(newUser.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Failed",
			"error": "Failed to hash password",
		})
        return
    }
    newUser.Password = hashedPassword

    query := "INSERT INTO users (id, username, email, password) VALUES (?, ?, ?, ?)"
    _, err = db.Exec(query, newUser.ID, newUser.Username, newUser.Email, newUser.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Failed",
			"error": "Failed to create user",
		})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"status": "success",
	})
}


func LoginUserHandler(c *gin.Context, db *sql.DB) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Failed",
			"error": "Invalid input",
		})
		return
	}

	var user User
	query := "SELECT id, username, password FROM users WHERE email = ?"
	err := db.QueryRow(query, loginData.Email).Scan(&user.ID, &user.Username, &user.Password)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "Failed",
			"error": "Invalid credentials",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Failed",
			"error": "Failed to query user",
		})
		return
	}

	if err := utils.CheckPassword(user.Password, loginData.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "Failed",
			"error": "Invalid credentials",
		})
		return
	}

	accessToken, refreshToken, err := utils.GenerateJWT(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Failed",
			"error": "Failed to generate tokens",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "Login successfully",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}