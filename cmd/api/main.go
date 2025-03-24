package main

import (
	"log"
	"swift-codes-api/internal/app"
	"swift-codes-api/internal/config"
	"swift-codes-api/internal/db"
)

func main() {
	cfg := config.LoadConfig()

	database, err := db.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	application := app.NewApp(database)
	defer application.Close()

	err = db.RunMigrations(database, "migrations")
	if err != nil {
		log.Fatalf("Migration error: %v", err)
	}

	log.Println("App initialized successfully!")
}
