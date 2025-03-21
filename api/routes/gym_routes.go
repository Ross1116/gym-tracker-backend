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
		gym.GET("/id/:id", func(c *gin.Context) {
			handlers.HandleGetGymByID(db, c)
		})
		gym.GET("/user/:user_id", func(c *gin.Context) {
			handlers.HandleGetGymsByUserID(db, c)
		})
		gym.PUT("/:id", func(c *gin.Context) {
			handlers.HandleUpdateGym(db, c)
		})
		gym.DELETE("/:id", func(c *gin.Context) {
			handlers.HandleDeleteGym(db, c)
		})
	}
}
