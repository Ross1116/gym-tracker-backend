package routes

import (
	"database/sql"

	"github.com/Ross1116/gym-tracker-backend/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupEquipmentTypeRoutes(db *sql.DB, router *gin.Engine) {
	equipmentTypes := router.Group("/api/equipment-types")
	{
		equipmentTypes.GET("", func(c *gin.Context) {
			handlers.HandleGetAllEquipmentTypes(db, c)
		})
		equipmentTypes.POST("", func(c *gin.Context) {
			handlers.HandleCreateEquipmentType(db, c)
		})
		equipmentTypes.GET("/:id", func(c *gin.Context) {
			handlers.HandleGetEquipmentType(db, c)
		})
		equipmentTypes.PUT("/:id", func(c *gin.Context) {
			handlers.HandleUpdateEquipmentType(db, c)
		})
		equipmentTypes.DELETE("/:id", func(c *gin.Context) {
			handlers.HandleDeleteEquipmentType(db, c)
		})
	}
}
