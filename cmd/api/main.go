package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"swift-codes-api/internal/config"
	"swift-codes-api/internal/db"
	"swift-codes-api/internal/handler"
	"swift-codes-api/internal/importer"
	"swift-codes-api/internal/repository"
	"swift-codes-api/internal/service"
)

func main() {
	cfg := config.LoadConfig()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pingCtx, pingCancel := context.WithTimeout(ctx, 5*time.Second)
	defer pingCancel()

	database, err := db.NewPostgresConnection(pingCtx, cfg)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	if err := db.RunMigrations(database, "migrations"); err != nil {
		log.Fatalf("Migration error: %v", err)
	}

	swiftRepo := repository.NewSwiftRepository(database)
	swiftService := service.NewSwiftService(swiftRepo)
	swiftHandler := handler.NewSwiftHandler(swiftService)

	importFile := "swift_data.xlsx"
	if err := importer.ImportSwiftCodesFromXLSX(ctx, importFile, swiftService); err != nil {
		log.Printf("IMPORT ERROR: %v", err)
	}

	router := chi.NewRouter()
	router.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		middleware.Timeout(15*time.Second),
	)
	router.Get("/v1/swift-codes/{swiftCode}", swiftHandler.GetSwiftCode)
	router.Get("/v1/swift-codes/country/{countryISO2}", swiftHandler.GetSwiftCodesByCountry)
	router.Post("/v1/swift-codes", swiftHandler.CreateSwiftCode)
	router.Delete("/v1/swift-codes/{swiftCode}", swiftHandler.DeleteSwiftCode)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		log.Println("Starting HTTP server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	log.Println("Shutting down HTTP server...")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP graceful shutdown failed: %v", err)
	}
	log.Println("Server stopped cleanly")
}
