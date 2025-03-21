package routes

import (
	"database/sql"

	"github.com/Ross1116/gym-tracker-backend/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupGymRoutes(db *sql.DB, router *gin.Engine) {
	gym := router.Group("/api/gyms")
	{
		gym.GET("", func(c *gin.Context) {
			handlers.HandleGetGyms(db, c)
		})
		gym.POST("", func(c *gin.Context) {
			handlers.HandleCreateGym(db, c)
		})
		gym.GET("/:id", func(c *gin.Context) {
			handlers.HandleGetGymByID(db, c)
		})
		gym.PUT("/:id")
		gym.DELETE("/:id")
	}
}
