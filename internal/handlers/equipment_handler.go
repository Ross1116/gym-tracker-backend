package handlers

import (
	"database/sql"
	"net/http"

	"github.com/Ross1116/gym-tracker-backend/internal/models"
	"github.com/gin-gonic/gin"
)

func HandleGetAllGymEquipments(db *sql.DB, c *gin.Context) {
	gymID := c.Param("gymId")

	query := `
			SELECT 
					ge.id, 
					ge.gym_id, 
					ge.equipment_type_id, 
					et.name AS equipment_name,
					ge.weight,
					ge.notes
			FROM gym_equipment ge
			JOIN equipment_types et ON ge.equipment_type_id = et.id
			WHERE ge.gym_id = $1
	`
	rows, err := db.Query(query, gymID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var equipments []models.GymEquipmentWithDetails
	for rows.Next() {
		var equipment models.GymEquipmentWithDetails
		if err := rows.Scan(
			&equipment.ID,
			&equipment.GymID,
			&equipment.EquipmentTypeID,
			&equipment.EquipmentName,
			&equipment.Weight,
			&equipment.Notes,
		); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		equipments = append(equipments, equipment)
	}

	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(equipments) == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No equipments found for this gym"})
		return
	}

	c.IndentedJSON(http.StatusOK, equipments)
}

func HandleAddNewGymEquipment(db *sql.DB, c *gin.Context) {}

func HandleGetGymEquipment(db *sql.DB, c *gin.Context) {}

func HandleUpdateGymEquipment(db *sql.DB, c *gin.Context) {}

func HandleDeleteGymEquipment(db *sql.DB, c *gin.Context) {}
