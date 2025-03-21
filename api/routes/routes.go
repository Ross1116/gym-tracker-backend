package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(db *sql.DB) {
	router := gin.Default()

	SetupGymRoutes(db, router)
	SetupEquipmentRoutes(db, router)

	router.Run(":8080")
}
