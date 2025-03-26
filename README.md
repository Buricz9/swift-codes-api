# Swift Codes API

A REST API project for managing banks' SWIFT codes, built in Go with PostgreSQL and containerized using Docker. The API allows retrieving bank information, searching SWIFT codes by country, adding new records, and deleting them. Initial data is imported from an Excel file (swift_data.xlsx) when the application starts.

## Technologies

- Go 1.24
- PostgreSQL 16
- Docker + Docker Compose
- Clean Architecture (handler → service → repository)
- Unit and integration tests

## Prerequisites
- Docker Desktop **(must be running before executing `docker-compose up`)**
- `docker-compose`
  
## How to Run the Project

1. **Clone the repository:**

   ```bash
   git clone https://github.com/Buricz9/swift-codes-api.git
   cd swift-codes-api
    ```
2. **Make sure Docker Desktop is running**

3. **Start the application and the database:**
   ```bash
   docker-compose up --build
   ```
  Note: Although go.mod and go.sum are included in the repository, if an error occurs when running via Docker (step 3), it's recommended to delete both files and run:
      ```bash
      go mod init swift-codes-api
      go mod tidy
      ```
  
The application will be available at: http://localhost:8080

## Data Import
Upon application startup, data is automatically imported from the swift_data.xlsx file. The file is located in the root directory of the project. The Dockerfile automatically copies this file into the container.

## Tests
Unit tests (with mocks):
```bash
go test ./internal/service -v
```
Integration tests (application + database running):
```bash
go test -tags=integration ./internal/integration -v
```
Before running integration tests, make sure the application is running (docker-compose up) and listening on localhost:8080.
Integration tests automatically clear the database before execution.

## Example Endpoints
```bash
GET /v1/swift-codes/BPKOPLPWXXX – pobierz dane HQ (z branchami)
GET /v1/swift-codes/BPKOPLPWXYZ – pobierz dane branch
GET /v1/swift-codes/country/PL – wszystkie SWIFTy z Polski
POST /v1/swift-codes – dodaj nowy kod SWIFT
DELETE /v1/swift-codes/{code} – usuń kod SWIFT
```

## Feature, nie bug
If you attempt to add a SWIFT code that already exists in the database, the application will notify you with the following messages:
```bash
Updated existing swift_code=   (w konsoli)
{"message":"Swift Code created successfully"} (przykładowo w postmana)
```
This behavior is intentional and ensures the uniqueness of SWIFT codes. The application does not allow duplicate entries – it always overwrites the previous version.
If you want to test adding the same code again, first delete the existing entry before sending a POST request.
