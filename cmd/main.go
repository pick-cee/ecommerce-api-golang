package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/pick-cee/go-ecommerce-api/internal/env"
)

func main() {
	ctx:= context.Background()
	cfg := config{
		addrr: ":8080",
		db: dbConfig{
			dsn: env.GetString("GOOSE_DBSTRING","host=localhost port=5432 user=postgres password=password dbname=ecommerce-go sslmode=disable"),
		},
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// databse
	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		panic(err)
	}

	logger.Info("Connected to database", "dsn", cfg.db.dsn)
	defer conn.Close(ctx)

	api := application{
		config: cfg,
		db: conn,
	}

	if err := api.run(api.mount()); err != nil {
		slog.Error("Server failed to start", "error", err)
	}
	os.Exit(1)
}