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
# config.ini is gitignored so we can't copy it at build time. Ship the example
# baked in; the real config is mounted in via docker-compose at runtime, or
# secrets can be provided via JWT_SECRET / DB_* env vars (see config.go).
COPY --from=builder /app/config.example.ini ./config.ini

EXPOSE 8080

ENTRYPOINT ["./saas-task-backend"]
