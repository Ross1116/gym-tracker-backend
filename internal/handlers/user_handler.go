package handlers

import (
	"database/sql"
	"net/http"

	"github.com/Ross1116/gym-tracker-backend/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func HandleGetUsers(db *sql.DB, c *gin.Context) {
	rows, err := db.Query("SELECT id, email, created_at, updated_at FROM users")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error fetching data"})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, u)
	}

	c.JSON(http.StatusOK, users)
}

func HandleCreateUser(db *sql.DB, c *gin.Context) {
	var user models.User

	if err := c.BindJSON(&user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
	}

	query := "INSERT INTO users (email, password_hash, created_at, updated_at) VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id;"
	err = db.QueryRow(query, user.Email, string(hashedPassword)).Scan(&user.ID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.PasswordHash = ""

	c.JSON(http.StatusCreated, user)
}
