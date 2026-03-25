# Build stage
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY cmd ./cmd
COPY internal ./internal
COPY tools ./tools

RUN swag init -g cmd/taxhub/main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /taxhub ./cmd/taxhub

# Runtime stage
FROM alpine:3.20

WORKDIR /app

RUN addgroup -S app && adduser -S app -G app

COPY --from=builder /taxhub /app/taxhub
COPY scientific_names.csv /app/scientific_names.csv

ENV ADDR=:8080
ENV CSV_PATH=/app/scientific_names.csv

EXPOSE 8080

USER app

CMD ["/app/taxhub"]
