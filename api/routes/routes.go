package routes

import (
	"database/sql"

	"github.com/Ross1116/gym-tracker-backend/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(db *sql.DB) {
	docs.SwaggerInfo.Host = "localhost:9000"

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Your Next.js app URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	SetupGymRoutes(db, router)
	SetupEquipmentRoutes(db, router)
	SetupUserRoutes(db, router)
	SetupEquipmentTypeRoutes(db, router)
	SetupExerciseRoutes(db, router)
	SetupWorkoutRoutes(db, router)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":9000")
}
