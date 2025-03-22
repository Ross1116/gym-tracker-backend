package main

import (
	"database/sql"
	"log"

	"github.com/Ross1116/gym-tracker-backend/api/routes"
	_ "github.com/lib/pq"
)

var db *sql.DB

// @title Gym Tracker API
// @version 1.0
// @description API for tracking gym workouts and exercises
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:9000
// @BasePath /api/
// @schemes http
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
}
