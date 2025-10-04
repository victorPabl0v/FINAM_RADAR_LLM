package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool() *pgxpool.Pool {
	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		dsn = "postgres://news_user:news_user@localhost:5432/news_db?sslmode=disable"
	}

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("unable to parse database URL: %v", err)
	}
	cfg.MaxConns = 50                        // максимум 50 подключений
	cfg.MinConns = 5                         // минимум 5 (держит активными)
	cfg.MaxConnLifetime = time.Hour          // пересоздаёт соединение раз в час
	cfg.MaxConnIdleTime = 10 * time.Minute   // разрывает неиспользуемые соединения
	cfg.HealthCheckPeriod = 30 * time.Second // проверяет соединения каждые 30 сек

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("unable to ping database: %v", err)
	}

	log.Println("database connected to", dsn)
	return pool
}
