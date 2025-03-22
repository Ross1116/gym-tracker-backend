package routes

import (
	"database/sql"

	"github.com/Ross1116/gym-tracker-backend/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(db *sql.DB, router *gin.Engine) {
	users := router.Group("/api/users")
	{
		users.GET("", func(c *gin.Context) {
			handlers.HandleGetUsers(db, c)
		})
		users.POST("", func(c *gin.Context) {
			handlers.HandleCreateUser(db, c)
		})
	}
}
