package integration

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

const baseURL = "http://localhost:8080/v1/swift-codes"

var db *sql.DB

func TestMain(m *testing.M) {
	var err error
	db, err = sql.Open("postgres", "host=localhost port=5432 user=swiftuser password=swiftpass dbname=swiftcodesdb sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	time.Sleep(3 * time.Second)

	clearDB()

	os.Exit(m.Run())
}

func clearDB() {
	db.Exec(`DELETE FROM swift_codes;`)
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

	resp, err := http.Post(baseURL, "application/json", bytes.NewBuffer(body))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestGetSwiftCode_HQ(t *testing.T) {
	resp, err := http.Get(baseURL + "/TESTCODE123")
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
	resp, err := http.Get(baseURL + "/country/PL")
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
	req, err := http.NewRequest(http.MethodDelete, baseURL+"/TESTCODE123", nil)
	assert.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetSwiftCode_NotFound(t *testing.T) {
	resp, err := http.Get(baseURL + "/DOESNOTEXIST")
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
