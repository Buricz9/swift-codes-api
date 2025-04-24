package integration

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	_ "net/http/httptest"
	"os"
	"testing"
	"time"

	_ "swift-codes-api/internal/app"
	"swift-codes-api/internal/config"
	"swift-codes-api/internal/db"
	"swift-codes-api/internal/handler"
	"swift-codes-api/internal/repository"
	"swift-codes-api/internal/service"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var (
	srv    *http.Server
	dbConn *sql.DB
)

func TestMain(m *testing.M) {
	// Load config and DB
	cfg := config.LoadConfig()
	var err error
	dbConn, err = db.NewPostgresConnection(db.Config(cfg))
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer dbConn.Close()

	// Run migrations
	if err := db.RunMigrations(dbConn, "migrations"); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Setup router and server
	repo := repository.NewSwiftRepository(dbConn)
	svc := service.NewSwiftService(repo)
	h := handler.NewSwiftHandler(svc)
	r := chi.NewRouter()
	r.Get("/v1/swift-codes/{swiftCode}", h.GetSwiftCode)
	r.Get("/v1/swift-codes/country/{countryISO2}", h.GetSwiftCodesByCountry)
	r.Post("/v1/swift-codes", h.CreateSwiftCode)
	r.Delete("/v1/swift-codes/{swiftCode}", h.DeleteSwiftCode)

	srv = &http.Server{Addr: ":8080", Handler: r}
	// Start server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for server to start
	time.Sleep(2 * time.Second)

	// Clear table
	clearDB()

	// Run tests
	code := m.Run()

	// Shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	os.Exit(code)
}

func clearDB() {
	// Delete all swift codes
	dbConn.ExecContext(context.Background(), "DELETE FROM swift.swift_codes;")
}

func TestCreateSwiftCode(t *testing.T) {
	payload := map[string]interface{}{
		"swiftCode":     "TESTCODE123",
		"bankName":      "Test Bank",
		"address":       "Test Address",
		"countryISO2":   "PL",
		"countryName":   "Poland",
		"isHeadquarter": true,
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post("http://localhost:8080/v1/swift-codes", "application/json", bytes.NewBuffer(body))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestGetSwiftCode_HQ(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/v1/swift-codes/TESTCODE123")
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	parseJSON(t, resp.Body, &result)

	assert.Equal(t, "TESTCODE123", result["swiftCode"])
	assert.Equal(t, true, result["isHeadquarter"])
	assert.Contains(t, result, "branches")
}

func TestGetSwiftCodesByCountry(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/v1/swift-codes/country/PL")
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	parseJSON(t, resp.Body, &result)

	assert.Equal(t, "PL", result["countryISO2"])
	assert.Contains(t, result, "swiftCodes")
	codes := result["swiftCodes"].([]interface{})
	assert.GreaterOrEqual(t, len(codes), 1)
}

func TestDeleteSwiftCode(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/v1/swift-codes/TESTCODE123", nil)
	assert.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetSwiftCode_NotFound(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/v1/swift-codes/DOESNOTEXIST")
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func parseJSON(t *testing.T, body io.Reader, target interface{}) {
	t.Helper()
	data, err := io.ReadAll(body)
	assert.NoError(t, err)
	err = json.Unmarshal(data, target)
	assert.NoError(t, err)
}
