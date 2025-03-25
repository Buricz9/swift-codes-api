# Swift Codes API

Projekt REST API do zarządzania kodami SWIFT banków, stworzony w języku Go z wykorzystaniem PostgreSQL i konteneryzacji za pomocą Dockera. API umożliwia pobieranie informacji o bankach, wyszukiwanie kodów SWIFT po kraju, dodawanie nowych rekordów, a także ich usuwanie. Dane początkowe są importowane z pliku Excel (`swift_data.xlsx`) podczas uruchamiania aplikacji.

## Technologie

- Go 1.24
- PostgreSQL 16
- Docker + Docker Compose
- Clean Architecture (handler → service → repository)
- Testy jednostkowe i integracyjne

## Wymagania wstępne
- Docker Desktop **(musi być uruchomiony przed komendą `docker-compose up`)**
- `docker-compose`

## Jak uruchomić projekt

1. **Sklonuj repozytorium:**

   ```bash
   git clone https://github.com/Buricz9/swift-codes-api.git
   cd swift-codes-api
    ```
2. **Upewnij się, że Docker Desktop jest uruchomiony.**

3. **Uruchom aplikację i bazę danych:**
   ```bash
   docker-compose up --build
   ```
Aplikacja będzie dostępna pod adresem: http://localhost:8080

## Import danych
Podczas startu aplikacji następuje automatyczny import danych z pliku swift_data.xlsx. Plik znajduje się w katalogu głównym projektu. Dockerfile automatycznie kopiuje ten plik do kontenera.

## Testy
Testy jednostkowe (mocki):
```bash
go test ./internal/service -v
```
Testy integracyjne (uruchomiona aplikacja + baza danych):
```bash
go test -tags=integration ./internal/integration -v
```
Przed uruchomieniem testów integracyjnych upewnij się, że aplikacja działa (docker-compose up) i nasłuchuje na localhost:8080.
Testy jednostkowe domyśnie opróżniają baze danych przed wykonaniem się

## Przykładowe endpointy
```bash
GET /v1/swift-codes/BPKOPLPWXXX – pobierz dane HQ (z branchami)
GET /v1/swift-codes/BPKOPLPWXYZ – pobierz dane branch
GET /v1/swift-codes/country/PL – wszystkie SWIFTy z Polski
POST /v1/swift-codes – dodaj nowy kod SWIFT
DELETE /v1/swift-codes/{code} – usuń kod SWIFT
```

## Feature, nie bug
W przypadku próby dodania wpisu z kodem SWIFT, który już istnieje w bazie danych, aplikacja poinformuje o tym komunikatem:
```bash
Updated existing swift_code=   (w konsoli)
{"message":"Swift Code created successfully"} (przykładowo w postmana)
```
Zachowanie to jest w pełni zamierzone i ma na celu zapewnienie unikalności kodów SWIFT. Aplikacja nie dopuszcza do tworzenia duplikatów - natomiast zawsze nadpisuje poprzendnie wrsje.
Jeśli chcesz przetestować dodanie kodu ponownie, usuń istniejący wpis przed wysłaniem żądania `POST`.
