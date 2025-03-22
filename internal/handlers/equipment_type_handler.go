package handlers

import (
	"database/sql"
	"net/http"

	"github.com/Ross1116/gym-tracker-backend/internal/models"
	"github.com/gin-gonic/gin"
)

func HandleCreateEquipmentType(db *sql.DB, c *gin.Context) {
	var input models.EquipmentTypeInput
	if err := c.BindJSON(&input); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM equipment_types WHERE name = $1)", input.Name).Scan(&exists)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if exists {
		c.IndentedJSON(http.StatusConflict, gin.H{"error": "Equipment type with this name already exists"})
		return
	}

	query := "INSERT INTO equipment_types (name) VALUES ($1) RETURNING id"

	var id int
	err = db.QueryRow(query, input.Name).Scan(&id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newEquipmentType := models.EquipmentType{
		ID:   id,
		Name: input.Name,
	}

	c.IndentedJSON(http.StatusCreated, newEquipmentType)
}
