package handlers

import (
	"database/sql"
	"net/http"

	"github.com/Ross1116/gym-tracker-backend/internal/models"
	"github.com/gin-gonic/gin"
)

func HandleGetGyms(db *sql.DB, c *gin.Context) {
	rows, err := db.Query("SELECT id, user_id, name, created_at FROM gyms")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error fetching data"})
		return
	}
	defer rows.Close()

	var gyms []models.Gym
	for rows.Next() {
		var gym models.Gym
		if err := rows.Scan(&gym.ID, &gym.UserID, &gym.Name, &gym.CreatedAt); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		gyms = append(gyms, gym)
	}

	c.JSON(http.StatusOK, gyms)
}

func HandleCreateGym(db *sql.DB, c *gin.Context) {
	var gym models.Gym
	if err := c.BindJSON(&gym); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if gym.Name == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Gym name is required"})
		return
	}

	query := "INSERT INTO gyms (user_id, name) VALUES ($1, $2) RETURNING id, created_at"
	err := db.QueryRow(query, gym.UserID, gym.Name).Scan(&gym.ID, &gym.CreatedAt)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create gym"})
		return
	}

	c.JSON(http.StatusCreated, gym)
}

func HandleGetGymByID(db *sql.DB, c *gin.Context) {
	id := c.Param("id")

	query := "SELECT id, user_id, name, created_at FROM gyms WHERE id=$1"
	row := db.QueryRow(query, id)

	var gym models.Gym
	if err := row.Scan(&gym.ID, &gym.UserID, &gym.Name, &gym.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Gym not found"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, gym)
}
