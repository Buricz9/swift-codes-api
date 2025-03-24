package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"swift-codes-api/internal/app"
	"swift-codes-api/internal/config"
	"swift-codes-api/internal/db"
	"swift-codes-api/internal/handler"
	"swift-codes-api/internal/importer"
	"swift-codes-api/internal/repository"
	"swift-codes-api/internal/service"
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

	swiftRepo := repository.NewSwiftRepository(database)
	swiftService := service.NewSwiftService(swiftRepo)
	swiftHandler := handler.NewSwiftHandler(swiftService)

	// Wywołanie importu z pliku XLSX:
	importFilePath := "swift_data.xlsx" // ścieżka do Twojego pliku
	if err := importer.ImportSwiftCodesFromXLSX(importFilePath, swiftService); err != nil {
		log.Printf("IMPORT ERROR: %v", err)
	}

	router := chi.NewRouter()
	router.Get("/v1/swift-codes/{swiftCode}", swiftHandler.GetSwiftCode)
	router.Get("/v1/swift-codes/country/{countryISO2}", swiftHandler.GetSwiftCodesByCountry)
	router.Post("/v1/swift-codes", swiftHandler.CreateSwiftCode)
	router.Delete("/v1/swift-codes/{swiftCode}", swiftHandler.DeleteSwiftCode)

	log.Println("Starting HTTP server on :8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}
