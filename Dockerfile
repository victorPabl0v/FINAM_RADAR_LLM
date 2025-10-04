# ---------- STAGE 1: Builder ----------
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Сборка бинарника
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./main.go

# ---------- STAGE 2: Runtime (с Go runtime) ----------
FROM golang:1.25-alpine AS runtime

# Устанавливаем сертификаты и tzdata (для HTTPS и корректного времени)
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Копируем бинарник и миграции из builder
COPY --from=builder /app/server .
COPY --from=builder /app/db/migrations ./db/migrations

# Переменные окружения
ENV DATABASE_URL=postgres://news_user:news_user@postgres:5432/news_db?sslmode=disable \
    PORT=8080 \
    GIN_MODE=release

EXPOSE 8080

ENTRYPOINT ["./server"]
