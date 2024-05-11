package main

import (
	"backend/pkg"
	"backend/pkg/db"
	"log"
	"net/http"
)

func main() {
	// Initialize the database connection
	db.InitDB("pkg/db/database.db")
	defer db.CloseDB()

	// Run migrations
	if err := db.RunMigrations(); err != nil {
		log.Fatal("Error applying migrations:", err)
	}

	// set up CORS middle ware
	handler := pkg.SetupRouter()

	// Start the HTTP server
	if err := http.ListenAndServe(":8000", handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
