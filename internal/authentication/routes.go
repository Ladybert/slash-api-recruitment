package authentication

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, db *sql.DB) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", func(c *gin.Context) {
			RegisterUserHandler(c, db)
		})
		auth.POST("/login", func(c *gin.Context) {
			LoginUserHandler(c, db)
		})
	}
}
