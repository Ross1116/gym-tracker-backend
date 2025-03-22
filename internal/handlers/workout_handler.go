package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/Ross1116/gym-tracker-backend/internal/models"
	"github.com/gin-gonic/gin"
)

func HandleGetUserWorkouts(db *sql.DB, c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	query := `
			SELECT 
					ws.id, 
					ws.user_id, 
					ws.gym_id, 
					ws.created_at
			FROM workout_sessions ws
			WHERE ws.user_id = $1
			ORDER BY ws.created_at DESC
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var workouts []models.WorkoutSession
	for rows.Next() {
		var workout models.WorkoutSession
		if err := rows.Scan(
			&workout.ID,
			&workout.UserID,
			&workout.GymID,
			&workout.CreatedAt,
		); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		workouts = append(workouts, workout)
	}

	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, workouts)
}

func HandleCreateWorkout(db *sql.DB, c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var input models.WorkoutSessionInput
	if err := c.BindJSON(&input); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer tx.Rollback()

	var workoutID int
	err = tx.QueryRow(
		"INSERT INTO workout_sessions (user_id, gym_id) VALUES ($1, $2) RETURNING id, created_at",
		userIDInt, input.GymID,
	).Scan(&workoutID)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err = tx.Commit(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var createdWorkout models.WorkoutSession
	err = db.QueryRow(
		"SELECT id, user_id, gym_id, created_at FROM workout_sessions WHERE id = $1",
		workoutID,
	).Scan(
		&createdWorkout.ID,
		&createdWorkout.UserID,
		&createdWorkout.GymID,
		&createdWorkout.CreatedAt,
	)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, createdWorkout)
}

func HandleGetExerciseHistory(db *sql.DB, c *gin.Context) {
	exerciseID := c.Param("exercise_id")
	equipmentID := c.Param("equipment_id")
	userID := c.Query("user_id")

	if userID == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	query := `
			SELECT 
					we.id,
					we.workout_session_id,
					we.exercise_id,
					e.name AS exercise_name,
					we.gym_equipment_id,
					et.name AS equipment_name,
					we.weight,
					we.reps,
					we.sets,
					we.created_at
			FROM workout_exercises we
			JOIN exercises e ON we.exercise_id = e.id
			JOIN workout_sessions ws ON we.workout_session_id = ws.id
			JOIN gym_equipment ge ON we.gym_equipment_id = ge.id
			JOIN equipment_types et ON ge.equipment_type_id = et.id
			WHERE we.exercise_id = $1 
			AND we.gym_equipment_id = $2
			AND ws.user_id = $3
			ORDER BY we.created_at DESC
			LIMIT 10
	`

	rows, err := db.Query(query, exerciseID, equipmentID, userID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var history []models.WorkoutExerciseWithDetails
	for rows.Next() {
		var exercise models.WorkoutExerciseWithDetails
		if err := rows.Scan(
			&exercise.ID,
			&exercise.WorkoutSessionID,
			&exercise.ExerciseID,
			&exercise.ExerciseName,
			&exercise.GymEquipmentID,
			&exercise.EquipmentName,
			&exercise.Weight,
			&exercise.Reps,
			&exercise.Sets,
			&exercise.CreatedAt,
		); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		history = append(history, exercise)
	}

	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, history)
}

func HandleGetLatestExercise(db *sql.DB, c *gin.Context) {
	exerciseID := c.Param("exercise_id")
	equipmentID := c.Param("equipment_id")
	userID := c.Query("user_id")

	if userID == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	query := `
			SELECT 
					we.id,
					we.workout_session_id,
					we.exercise_id,
					e.name AS exercise_name,
					we.gym_equipment_id,
					et.name AS equipment_name,
					we.weight,
					we.reps,
					we.sets,
					we.created_at
			FROM workout_exercises we
			JOIN exercises e ON we.exercise_id = e.id
			JOIN workout_sessions ws ON we.workout_session_id = ws.id
			JOIN gym_equipment ge ON we.gym_equipment_id = ge.id
			JOIN equipment_types et ON ge.equipment_type_id = et.id
			WHERE we.exercise_id = $1 
			AND we.gym_equipment_id = $2
			AND ws.user_id = $3
			ORDER BY we.created_at DESC
			LIMIT 1
	`

	var exercise models.WorkoutExerciseWithDetails
	err := db.QueryRow(query, exerciseID, equipmentID, userID).Scan(
		&exercise.ID,
		&exercise.WorkoutSessionID,
		&exercise.ExerciseID,
		&exercise.ExerciseName,
		&exercise.GymEquipmentID,
		&exercise.EquipmentName,
		&exercise.Weight,
		&exercise.Reps,
		&exercise.Sets,
		&exercise.CreatedAt,
	)

	if err == sql.ErrNoRows {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No previous workout found for this exercise and equipment"})
		return
	} else if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, exercise)
}
