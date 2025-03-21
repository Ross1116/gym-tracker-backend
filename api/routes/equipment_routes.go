package routes

import (
	"database/sql"

	"github.com/Ross1116/gym-tracker-backend/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupEquipmentRoutes(db *sql.DB, router *gin.Engine) {
	gymEquipment := router.Group("/api/gyms/:gymId/equipment")
	{
		gymEquipment.GET("", func(c *gin.Context) {
			handlers.HandleGetAllGymEquipments(db, c)
		})

		gymEquipment.POST("", func(c *gin.Context) {
			handlers.HandleAddNewGymEquipment(db, c)
		})
	}

	equipmentRoutes := router.Group("/api/gym-equipment")
	{
		equipmentRoutes.POST("", func(c *gin.Context) {
			handlers.HandleGetGymEquipment(db, c)
		})

		equipmentRoutes.PUT("/:id", func(c *gin.Context) {
			handlers.HandleUpdateGymEquipment(db, c)
		})

		equipmentRoutes.DELETE("/:id", func(c *gin.Context) {
			handlers.HandleDeleteGymEquipment(db, c)
		})
	}
}
