package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/Ross1116/gym-tracker-backend/internal/models"
	"github.com/gin-gonic/gin"
)

func HandleGetAllExercises(db *sql.DB, c *gin.Context) {
	rows, err := db.Query("SELECT id, name FROM exercises ORDER BY name")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var exercises []models.Exercise
	for rows.Next() {
		var exercise models.Exercise
		if err := rows.Scan(&exercise.ID, &exercise.Name); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		exercises = append(exercises, exercise)
	}

	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, exercises)
}

func HandleCreateExercise(db *sql.DB, c *gin.Context) {
	var input models.ExerciseInput
	if err := c.BindJSON(&input); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM exercises WHERE name = $1)", input.Name).Scan(&exists)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if exists {
		c.IndentedJSON(http.StatusConflict, gin.H{"error": "Exercise with this name already exists"})
		return
	}

	var id int
	err = db.QueryRow("INSERT INTO exercises (name) VALUES ($1) RETURNING id", input.Name).Scan(&id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newExercise := models.Exercise{
		ID:   id,
		Name: input.Name,
	}

	c.IndentedJSON(http.StatusCreated, newExercise)
}

func HandleUpdateExercise(db *sql.DB, c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var input models.ExerciseInput
	if err := c.BindJSON(&input); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM exercises WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !exists {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
		return
	}

	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM exercises WHERE name = $1 AND id != $2)",
		input.Name, id).Scan(&exists)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if exists {
		c.IndentedJSON(http.StatusConflict, gin.H{"error": "Exercise with this name already exists"})
		return
	}

	_, err = db.Exec("UPDATE exercises SET name = $1 WHERE id = $2", input.Name, id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updatedExercise := models.Exercise{
		ID:   id,
		Name: input.Name,
	}

	c.IndentedJSON(http.StatusOK, updatedExercise)
}

func HandleDeleteExercise(db *sql.DB, c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var inUse bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM workout_exercises WHERE exercise_id = $1)", id).Scan(&inUse)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if inUse {
		c.IndentedJSON(http.StatusConflict, gin.H{"error": "Cannot delete exercise that is used in workouts"})
		return
	}

	result, err := db.Exec("DELETE FROM exercises WHERE id = $1", id)
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
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Exercise deleted successfully"})
}
