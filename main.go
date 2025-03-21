package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/Ross1116/gym-tracker-backend/api/routes"
	_ "github.com/lib/pq"
)

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Gym struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Exercise struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Equipment struct {
	ID    int    `json:"id"`
	GymID int    `json:"gym_id"`
	Name  string `json:"name"`
}

type WorkoutSession struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	ExerciseID  int       `json:"exercise_id"`
	EquipmentID int       `json:"equipment_id"`
	Weight      float64   `json:"weight"`
	Reps        int       `json:"reps"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PantryItem struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	Name            string    `json:"name"`
	Quantity        float64   `json:"quantity"`
	Unit            string    `json:"unit"`
	Threshold       float64   `json:"threshold"`
	CaloriesPerUnit float64   `json:"calories_per_unit"`
	ProteinPerUnit  float64   `json:"protein_per_unit"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Meal struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type MealIngredient struct {
	MealID       int     `json:"meal_id"`
	PantryItemID int     `json:"pantry_item_id"`
	QuantityUsed float64 `json:"quantity_used"`
}

type ShoppingList struct {
	UserID         int     `json:"user_id"`
	Name           string  `json:"name"`
	QuantityNeeded float64 `json:"quantity_needed"`
	Unit           string  `json:"unit"`
}

var db *sql.DB

func main() {
	var err error
	connStr := "host=localhost port=5432 user=admin password=admin dbname=mydb sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Could not connect to the db", err)
	}

	routes.SetupRoutes(db)

	// router := gin.Default()
	// router.GET("/users", handleGetUsers)
	// router.POST("/users", handleCreateUser)
	// router.Run(":8080")
}

// func handleCreateUser(c *gin.Context) {
// 	var user User

// 	if err := c.BindJSON(&user); err != nil {
// 		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
// 	if err != nil {
// 		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
// 	}

// 	query := "INSERT INTO users (email, password_hash, created_at, updated_at) VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id;"
// 	err = db.QueryRow(query, user.Email, string(hashedPassword)).Scan(&user.ID)
// 	if err != nil {
// 		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	user.PasswordHash = ""

// 	c.JSON(http.StatusCreated, user)
// }

// func handleGetUsers(c *gin.Context) {
// 	rows, err := db.Query("SELECT id, email, created_at, updated_at FROM users")
// 	if err != nil {
// 		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error fetching data"})
// 		return
// 	}
// 	defer rows.Close()

// 	var users []User
// 	for rows.Next() {
// 		var u User
// 		if err := rows.Scan(&u.ID, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
// 			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 		users = append(users, u)
// 	}

// 	c.JSON(http.StatusOK, users)
// }
