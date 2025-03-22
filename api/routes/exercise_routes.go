package routes

import (
	"database/sql"

	"github.com/Ross1116/gym-tracker-backend/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupExerciseRoutes(db *sql.DB, router *gin.Engine) {
	exercises := router.Group("/api/exercises")
	{
		exercises.GET("", func(c *gin.Context) {
			handlers.HandleGetAllExercises(db, c)
		})

		exercises.POST("", func(c *gin.Context) {
			handlers.HandleCreateExercise(db, c)
		})

		exercises.PUT("/:id", func(c *gin.Context) {
			handlers.HandleUpdateExercise(db, c)
		})

		exercises.DELETE("/:id", func(c *gin.Context) {
			handlers.HandleDeleteExercise(db, c)
		})
	}
}
