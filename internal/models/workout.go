package models

import (
	"time"
)

type WorkoutSession struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	GymID     int       `json:"gym_id"`
	CreatedAt time.Time `json:"created_at"`
}

type WorkoutSessionInput struct {
	GymID int `json:"gym_id" binding:"required"`
}

type WorkoutSessionWithExercises struct {
	WorkoutSession
	Exercises []WorkoutExerciseWithDetails `json:"exercises"`
}

type WorkoutExercise struct {
	ID               int       `json:"id"`
	WorkoutSessionID int       `json:"workout_session_id"`
	ExerciseID       int       `json:"exercise_id"`
	GymEquipmentID   int       `json:"gym_equipment_id"`
	Weight           float64   `json:"weight"`
	Reps             int       `json:"reps"`
	Sets             int       `json:"sets"`
	CreatedAt        time.Time `json:"created_at"`
}

type WorkoutExerciseInput struct {
	ExerciseID     int     `json:"exercise_id" binding:"required"`
	GymEquipmentID int     `json:"gym_equipment_id" binding:"required"`
	Weight         float64 `json:"weight" binding:"required"`
	Reps           int     `json:"reps" binding:"required"`
	Sets           int     `json:"sets" binding:"required"`
}

type WorkoutExerciseWithDetails struct {
	ID               int       `json:"id"`
	WorkoutSessionID int       `json:"workout_session_id"`
	ExerciseID       int       `json:"exercise_id"`
	ExerciseName     string    `json:"exercise_name"`
	GymEquipmentID   int       `json:"gym_equipment_id"`
	EquipmentName    string    `json:"equipment_name"`
	Weight           float64   `json:"weight"`
	Reps             int       `json:"reps"`
	Sets             int       `json:"sets"`
	CreatedAt        time.Time `json:"created_at"`
}

type WorkoutSessionWithExercisesInput struct {
	GymID     int                    `json:"gym_id" binding:"required"`
	Exercises []WorkoutExerciseInput `json:"exercises" binding:"required"`
}
