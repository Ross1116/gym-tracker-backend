package routes

import (
	"database/sql"

	"github.com/Ross1116/gym-tracker-backend/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupWorkoutRoutes(db *sql.DB, router *gin.Engine) {
	workouts := router.Group("/api/workouts")
	{
		workouts.GET("", func(c *gin.Context) {
			handlers.HandleGetUserWorkouts(db, c)
		})

		workouts.POST("", func(c *gin.Context) {
			handlers.HandleCreateWorkout(db, c)
		})

		workouts.GET("history/:exercise_id/:equipment_id", func(c *gin.Context) {
			handlers.HandleGetExerciseHistory(db, c)
		})

		workouts.GET("latest/:exercise_id/:equipment_id", func(c *gin.Context) {
			handlers.HandleGetLatestExercise(db, c)
		})
	}
}
