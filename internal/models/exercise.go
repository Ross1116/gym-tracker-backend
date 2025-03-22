package models

type Exercise struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ExerciseInput struct {
	Name string `json:"name" binding:"required"`
}
