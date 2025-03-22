package main

import (
	"database/sql"
	"log"

	"github.com/Ross1116/gym-tracker-backend/api/routes"
	_ "github.com/lib/pq"
)

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
}
