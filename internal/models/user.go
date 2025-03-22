package models

import "time"

type User struct {
	ID           int       `json:"id" swaggerignore:"true"`
	Email        string    `json:"email" example:"user@example.com"`
	PasswordHash string    `json:"password" example:"MySecurePassword123" swaggerignore:"true"`
	CreatedAt    time.Time `json:"created_at" swaggerignore:"true"`
	UpdatedAt    time.Time `json:"updated_at" swaggerignore:"true"`
}
