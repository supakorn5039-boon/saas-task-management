# Stage 1: Builder
FROM golang:1.25.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bin/saas-task-backend ./src/main.go

# Stage 2: Runner
FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/bin/saas-task-backend ./saas-task-backend
COPY --from=builder /app/config.ini ./config.ini

EXPOSE 8080

ENTRYPOINT ["./saas-task-backend"]
