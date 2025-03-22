package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/Ross1116/gym-tracker-backend/internal/models"
	"github.com/gin-gonic/gin"
)

// HandleGetAllGymEquipments godoc
// @Summary Get all equipment for a specific gym
// @Description Retrieve all equipment items associated with a gym
// @Tags GymEquipment
// @Accept json
// @Produce json
// @Param gymId path int true "ID of the gym"
// @Success 200 {array} models.GymEquipmentWithDetails
// @Failure 404 {object} models.ErrorResponse "No equipments found for this gym"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /gyms/{gymId}/equipment [get]
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

// HandleAddNewGymEquipment godoc
// @Summary Add new equipment to a gym
// @Description Add a new equipment item to a specific gym
// @Tags GymEquipment
// @Accept json
// @Produce json
// @Param gymId path int true "ID of the gym"
// @Param equipment body models.GymEquipmentInput true "Equipment details"
// @Success 201 {object} models.GymEquipment
// @Failure 400 {object} models.ErrorResponse "Invalid gym ID format or invalid input"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /gyms/{gymId}/equipment [post]
func HandleAddNewGymEquipment(db *sql.DB, c *gin.Context) {
	gymId := c.Param("gymId")

	gymIDInt, err := strconv.Atoi(gymId)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid gym ID format"})
		return
	}

	var input models.GymEquipmentInput
	if err := c.BindJSON(&input); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
                INSERT INTO gym_equipment (gym_id, equipment_type_id, weight, notes)
                VALUES ($1, $2, $3, $4)
                RETURNING id
        `

	var id int
	err = db.QueryRow(
		query,
		gymIDInt,
		input.EquipmentTypeID,
		input.Weight,
		input.Notes,
	).Scan(&id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newEquipment := models.GymEquipment{
		ID:              id,
		GymID:           gymIDInt,
		EquipmentTypeID: input.EquipmentTypeID,
		Weight:          input.Weight,
		Notes:           input.Notes,
	}

	c.IndentedJSON(http.StatusCreated, newEquipment)
}

// HandleGetGymEquipment godoc
// @Summary Get specific gym equipment by ID
// @Description Retrieve details of a specific gym equipment item
// @Tags GymEquipment
// @Accept json
// @Produce json
// @Param id path int true "ID of the equipment"
// @Success 200 {object} models.GymEquipmentWithDetails
// @Failure 400 {object} models.ErrorResponse "Invalid equipment ID format"
// @Failure 404 {object} models.ErrorResponse "Equipment not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /equipment/{id} [get]
func HandleGetGymEquipment(db *sql.DB, c *gin.Context) {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid equipment ID format"})
		return
	}

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
			WHERE ge.id = $1
	`

	var equipment models.GymEquipmentWithDetails
	err = db.QueryRow(query, idInt).Scan(
		&equipment.ID,
		&equipment.GymID,
		&equipment.EquipmentTypeID,
		&equipment.EquipmentName,
		&equipment.Weight,
		&equipment.Notes,
	)

	if err == sql.ErrNoRows {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Equipment not found"})
		return
	} else if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, equipment)
}

// HandleUpdateGymEquipment godoc
// @Summary Update gym equipment
// @Description Update the details of an existing gym equipment item
// @Tags GymEquipment
// @Accept json
// @Produce json
// @Param id path int true "ID of the equipment to update"
// @Param equipment body models.GymEquipmentInput true "Updated equipment details"
// @Success 200 {object} models.GymEquipmentWithDetails
// @Failure 400 {object} models.ErrorResponse "Invalid equipment ID format, invalid input, or equipment type not found"
// @Failure 404 {object} models.ErrorResponse "Equipment not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /equipment/{id} [put]
func HandleUpdateGymEquipment(db *sql.DB, c *gin.Context) {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid equipment ID format"})
		return
	}

	var input models.GymEquipmentInput
	if err := c.BindJSON(&input); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM gym_equipment WHERE id = $1)", idInt).Scan(&exists)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !exists {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Equipment not found"})
		return
	}

	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM equipment_types WHERE id = $1)", input.EquipmentTypeID).Scan(&exists)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !exists {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Equipment type not found"})
		return
	}

	query := `
			UPDATE gym_equipment 
			SET equipment_type_id = $1, weight = $2, notes = $3
			WHERE id = $4
			RETURNING gym_id
	`

	var gymID int
	err = db.QueryRow(
		query,
		input.EquipmentTypeID,
		input.Weight,
		input.Notes,
		idInt,
	).Scan(&gymID)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	detailsQuery := `
			SELECT 
					ge.id, 
					ge.gym_id, 
					ge.equipment_type_id, 
					et.name AS equipment_name,
					ge.weight,
					ge.notes
			FROM gym_equipment ge
			JOIN equipment_types et ON ge.equipment_type_id = et.id
			WHERE ge.id = $1
	`

	var updatedEquipment models.GymEquipmentWithDetails
	err = db.QueryRow(detailsQuery, idInt).Scan(
		&updatedEquipment.ID,
		&updatedEquipment.GymID,
		&updatedEquipment.EquipmentTypeID,
		&updatedEquipment.EquipmentName,
		&updatedEquipment.Weight,
		&updatedEquipment.Notes,
	)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, updatedEquipment)
}

// HandleDeleteGymEquipment godoc
// @Summary Delete gym equipment
// @Description Delete an existing gym equipment item
// @Tags GymEquipment
// @Accept json
// @Produce json
// @Param id path int true "ID of the equipment to delete"
// @Success 200 {object} models.SuccessResponse "Equipment successfully deleted"
// @Failure 400 {object} models.ErrorResponse "Invalid equipment ID format"
// @Failure 404 {object} models.ErrorResponse "Equipment not found"
// @Failure 409 {object} models.ErrorResponse "Cannot delete equipment that is in use"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /equipment/{id} [delete]
func HandleDeleteGymEquipment(db *sql.DB, c *gin.Context) {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid equipment ID format"})
		return
	}

	var inUse bool
	err = db.QueryRow(`
			SELECT EXISTS(
					SELECT 1 FROM workout_equipment 
					WHERE gym_equipment_id = $1
			)`, idInt).Scan(&inUse)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if inUse {
		c.IndentedJSON(http.StatusConflict, gin.H{"error": "Cannot delete equipment that is used in workout sessions"})
		return
	}

	result, err := db.Exec("DELETE FROM gym_equipment WHERE id = $1", idInt)
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
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Equipment not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Equipment removed successfully"})
}
