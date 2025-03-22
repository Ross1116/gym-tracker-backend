package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

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

func HandleGetWorkoutWithExercises(db *sql.DB, c *gin.Context) {
	workoutIDStr := c.Param("id")
	workoutID, err := strconv.Atoi(workoutIDStr)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid workout ID format"})
		return
	}

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

	var workout models.WorkoutSession
	err = db.QueryRow(
		"SELECT id, user_id, gym_id, created_at FROM workout_sessions WHERE id = $1 AND user_id = $2",
		workoutID, userIDInt,
	).Scan(&workout.ID, &workout.UserID, &workout.GymID, &workout.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Workout not found or not authorized"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	rows, err := db.Query(`
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
        JOIN exercises e ON e.id = we.exercise_id
        JOIN gym_equipment ge ON ge.id = we.gym_equipment_id
        JOIN equipment_types et ON et.id = ge.equipment_type_id
        WHERE we.workout_session_id = $1
        ORDER BY we.id
    `, workoutID)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	exercises := []models.WorkoutExerciseWithDetails{}
	for rows.Next() {
		var exercise models.WorkoutExerciseWithDetails
		err := rows.Scan(
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
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		exercises = append(exercises, exercise)
	}

	if err = rows.Err(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := models.WorkoutSessionWithExercises{
		WorkoutSession: workout,
		Exercises:      exercises,
	}

	c.IndentedJSON(http.StatusOK, result)
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

	var sessionInput models.WorkoutSessionInput
	if err := c.BindJSON(&sessionInput); err != nil {
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
	var createdAt time.Time

	err = tx.QueryRow(
		"INSERT INTO workout_sessions (user_id, gym_id) VALUES ($1, $2) RETURNING id, created_at",
		userIDInt, sessionInput.GymID,
	).Scan(&workoutID, &createdAt)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err = tx.Commit(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	createdWorkout := models.WorkoutSession{
		ID:        workoutID,
		UserID:    userIDInt,
		GymID:     sessionInput.GymID,
		CreatedAt: createdAt,
	}

	c.IndentedJSON(http.StatusCreated, createdWorkout)
}

func HandleCreateWorkoutWithExercises(db *sql.DB, c *gin.Context) {
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

	var input models.WorkoutSessionWithExercisesInput
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
	var createdAt time.Time

	err = tx.QueryRow(
		"INSERT INTO workout_sessions (user_id, gym_id) VALUES ($1, $2) RETURNING id, created_at",
		userIDInt, input.GymID,
	).Scan(&workoutID, &createdAt)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	exerciseDetails := make([]models.WorkoutExerciseWithDetails, 0, len(input.Exercises))

	for _, exercise := range input.Exercises {
		var exerciseID int
		var exerciseCreatedAt time.Time
		var exerciseName, equipmentName string

		err = tx.QueryRow(
			`INSERT INTO workout_exercises 
            (workout_session_id, exercise_id, gym_equipment_id, weight, reps, sets) 
            VALUES ($1, $2, $3, $4, $5, $6) 
            RETURNING id, created_at`,
			workoutID,
			exercise.ExerciseID,
			exercise.GymEquipmentID,
			exercise.Weight,
			exercise.Reps,
			exercise.Sets,
		).Scan(&exerciseID, &exerciseCreatedAt)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to add exercise: " + err.Error()})
			return
		}

		err = tx.QueryRow(
			`SELECT e.name, et.name 
             FROM exercises e
             JOIN gym_equipment ge ON ge.id = $1
             JOIN equipment_types et ON et.id = ge.equipment_type_id
             WHERE e.id = $2`,
			exercise.GymEquipmentID, exercise.ExerciseID,
		).Scan(&exerciseName, &equipmentName)
		if err != nil {
			exerciseName = "Unknown"
			equipmentName = "Unknown"
		}

		exerciseDetails = append(exerciseDetails, models.WorkoutExerciseWithDetails{
			ID:               exerciseID,
			WorkoutSessionID: workoutID,
			ExerciseID:       exercise.ExerciseID,
			ExerciseName:     exerciseName,
			GymEquipmentID:   exercise.GymEquipmentID,
			EquipmentName:    equipmentName,
			Weight:           exercise.Weight,
			Reps:             exercise.Reps,
			Sets:             exercise.Sets,
			CreatedAt:        exerciseCreatedAt,
		})
	}

	if err = tx.Commit(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	createdWorkout := models.WorkoutSessionWithExercises{
		WorkoutSession: models.WorkoutSession{
			ID:        workoutID,
			UserID:    userIDInt,
			GymID:     input.GymID,
			CreatedAt: createdAt,
		},
		Exercises: exerciseDetails,
	}

	c.IndentedJSON(http.StatusCreated, createdWorkout)
}

func HandleAddWorkoutExercise(db *sql.DB, c *gin.Context) {
	sessionIDStr := c.Param("sessionId")
	sessionID, err := strconv.Atoi(sessionIDStr)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid workout session ID"})
		return
	}

	var exerciseInput models.WorkoutExerciseInput
	if err := c.BindJSON(&exerciseInput); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer tx.Rollback()

	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM workout_sessions WHERE id = $1)", sessionID).Scan(&exists)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Workout session not found"})
		return
	}

	var exerciseID int
	var createdAt time.Time
	err = tx.QueryRow(
		`INSERT INTO workout_exercises 
        (workout_session_id, exercise_id, gym_equipment_id, weight, reps, sets) 
        VALUES ($1, $2, $3, $4, $5, $6) 
        RETURNING id, created_at`,
		sessionID,
		exerciseInput.ExerciseID,
		exerciseInput.GymEquipmentID,
		exerciseInput.Weight,
		exerciseInput.Reps,
		exerciseInput.Sets,
	).Scan(&exerciseID, &createdAt)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err = tx.Commit(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	createdExercise := models.WorkoutExercise{
		ID:               exerciseID,
		WorkoutSessionID: sessionID,
		ExerciseID:       exerciseInput.ExerciseID,
		GymEquipmentID:   exerciseInput.GymEquipmentID,
		Weight:           exerciseInput.Weight,
		Reps:             exerciseInput.Reps,
		Sets:             exerciseInput.Sets,
		CreatedAt:        createdAt,
	}

	c.IndentedJSON(http.StatusCreated, createdExercise)
}
