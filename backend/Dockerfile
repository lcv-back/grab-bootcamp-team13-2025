# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go mod tidy
RUN go build -o isymptom ./cmd/main.go
RUN go build -o email_worker ./cmd/email_worker/main.go
RUN ls -lah /app
COPY .env .env

# Production stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
RUN mkdir -p /app/logs

COPY --from=builder /app/isymptom .
COPY .env .env
COPY --from=builder /app/email_worker .
RUN chmod +x /app/email_worker
COPY --from=builder /app/config/config.yaml ./config/
# Copy mock outbreak data for demo
COPY --from=builder /app/internal/infrastructure/mock_outbreaks.json ./internal/infrastructure/mock_outbreaks.json

ENV JWT_SECRET="" \
    DB_HOST="" \
    DB_PORT="" \
    DB_NAME="" \
    DB_USER="" \
    DB_PASSWORD="" \
    RABBITMQ_URL="" \
    REDIS_URL="" \
    EMAIL_HOST="" \
    EMAIL_PORT="" \
    EMAIL_USERNAME="" \
    EMAIL_PASSWORD="" \
    MINIO_ENDPOINT="" \
    MINIO_ACCESS_KEY="" \
    MINIO_SECRET_KEY="" \
    MINIO_BUCKET="" \
    MINIO_USE_SSL="" \
    ML_URL=""

EXPOSE 3001

CMD ["./isymptom"]