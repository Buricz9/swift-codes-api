package app

import (
	"database/sql"
	"log"
)

type App struct {
	DB *sql.DB
}

func NewApp(db *sql.DB) *App {
	return &App{
		DB: db,
	}
}

func (a *App) Close() {
	if err := a.DB.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}
}
