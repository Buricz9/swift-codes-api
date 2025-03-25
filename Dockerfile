FROM golang:1.24.0-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o swift-codes-api ./cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=builder /app/swift-codes-api .
COPY --from=builder /app/migrations ./migrations
COPY swift_data.xlsx .
EXPOSE 8080
ENTRYPOINT ["./swift-codes-api"]
