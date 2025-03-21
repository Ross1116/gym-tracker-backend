package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetupEquipmentTypeRoutes(db *sql.DB, router *gin.Engine) {
	equipmentTypes := router.Group("/api/equipment-types")
	{
		equipmentTypes.GET("")
		equipmentTypes.POST("")
		equipmentTypes.GET("/:id")
		equipmentTypes.PUT("/:id")
		equipmentTypes.DELETE("/:id")
	}
}
