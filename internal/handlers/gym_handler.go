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

func HandleGetGymsByUserID(db *sql.DB, c *gin.Context) {
	id := c.Param("user_id")

	query := "SELECT id, user_id, name, created_at FROM gyms WHERE user_id=$1"
	rows, err := db.Query(query, id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(gyms) == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No gyms found for this user"})
		return
	}

	c.IndentedJSON(http.StatusOK, gyms)
}

func HandleUpdateGym(db *sql.DB, c *gin.Context) {
	id := c.Param("id")

	var gym models.Gym
	if err := c.BindJSON(&gym); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	query := "UPDATE gyms SET user_id=$2, name=$3 WHERE id=$1 RETURNING id, user_id, name, created_at"
	err := db.QueryRow(query, id, gym.UserID, gym.Name).Scan(&gym.ID, &gym.UserID, &gym.Name, &gym.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Gym not found"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update gym"})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, gym)
}

func HandleDeleteGym(db *sql.DB, c *gin.Context) {
	id := c.Param("id")

	var gym models.Gym
	getQuery := "SELECT id, user_id, name, created_at FROM gyms WHERE id=$1"
	err := db.QueryRow(getQuery, id).Scan(&gym.ID, &gym.UserID, &gym.Name, &gym.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Gym not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	query := "DELETE FROM gyms WHERE id=$1"
	_, err = db.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete gym"})
		return
	}

	c.IndentedJSON(http.StatusOK, "Deleted gym sucessfully")
}
