package routes

import (
	"database/sql"

	"github.com/Ross1116/gym-tracker-backend/docs"
	_ "github.com/Ross1116/gym-tracker-backend/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(db *sql.DB) {
	docs.SwaggerInfo.Host = "localhost:9000"

	router := gin.Default()

	SetupGymRoutes(db, router)
	SetupEquipmentRoutes(db, router)
	SetupUserRoutes(db, router)
	SetupEquipmentTypeRoutes(db, router)
	SetupExerciseRoutes(db, router)
	SetupWorkoutRoutes(db, router)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":9000")
}
