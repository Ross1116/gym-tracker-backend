package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/Ross1116/gym-tracker-backend/internal/models"
	"github.com/gin-gonic/gin"
)

func HandleGetAllEquipmentTypes(db *sql.DB, c *gin.Context) {
	rows, err := db.Query("SELECT id, name FROM equipment_types ORDER BY name")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	equipmentTypes := []models.EquipmentType{}
	for rows.Next() {
		var equipmentType models.EquipmentType
		if err := rows.Scan(&equipmentType.ID, &equipmentType.Name); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		equipmentTypes = append(equipmentTypes, equipmentType)
	}

	c.IndentedJSON(http.StatusOK, equipmentTypes)
}

func HandleGetEquipmentType(db *sql.DB, c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var equipmentType models.EquipmentType
	err = db.QueryRow("SELECT id, name FROM equipment_types WHERE id = $1", id).Scan(
		&equipmentType.ID, &equipmentType.Name)

	if err == sql.ErrNoRows {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Equipment type not found"})
		return
	} else if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, equipmentType)
}

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

func HandleUpdateEquipmentType(db *sql.DB, c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var input models.EquipmentTypeInput
	if err := c.BindJSON(&input); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM equipment_types WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !exists {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Equipment type not found"})
		return
	}

	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM equipment_types WHERE name = $1 AND id != $2)",
		input.Name, id).Scan(&exists)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if exists {
		c.IndentedJSON(http.StatusConflict, gin.H{"error": "Equipment type with this name already exists"})
		return
	}

	_, err = db.Exec("UPDATE equipment_types SET name = $1 WHERE id = $2", input.Name, id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updatedEquipmentType := models.EquipmentType{
		ID:   id,
		Name: input.Name,
	}

	c.IndentedJSON(http.StatusOK, updatedEquipmentType)
}

func HandleDeleteEquipmentType(db *sql.DB, c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var inUse bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM gym_equipment WHERE equipment_type_id = $1)", id).Scan(&inUse)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if inUse {
		c.IndentedJSON(http.StatusConflict, gin.H{"error": "Cannot delete equipment type that is in use"})
		return
	}

	result, err := db.Exec("DELETE FROM equipment_types WHERE id = $1", id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Equipment type not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Equipment type deleted successfully"})
}
