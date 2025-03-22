package routes

import (
	"database/sql"

	"github.com/Ross1116/gym-tracker-backend/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupEquipmentTypeRoutes(db *sql.DB, router *gin.Engine) {
	equipmentTypes := router.Group("/api/equipment-types")
	{
		equipmentTypes.GET("")
		equipmentTypes.POST("", func(c *gin.Context) {
			handlers.HandleCreateEquipmentType(db, c)
		})
		equipmentTypes.GET("/:id")
		equipmentTypes.PUT("/:id")
		equipmentTypes.DELETE("/:id")
	}
}
