package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"swift-codes-api/internal/db"
)

func LoadConfig() db.Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env variables")
	}

	port, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		log.Fatalf("Invalid DB_PORT: %v", err)
	}

	return db.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     port,
		User:     getEnv("DB_USER", "swiftuser"),
		Password: getEnv("DB_PASSWORD", "swiftpass"),
		DBName:   getEnv("DB_NAME", "swiftcodesdb"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
