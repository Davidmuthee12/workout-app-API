package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Davidmuthee12/kicker/internals/env"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main () {
	godotenv.Load()
	ctx := context.Background()


	cfg := config{
		addr: ":" + env.GetString("API_PORT", "8080"),
		db: dbConfig{
			dsn: "postgres://" + env.GetString("DB_USER", "postgres") + ":" + env.GetString("DB_PASSWORD", "postgres") + "@" + env.GetString("DB_HOST", "localhost") + ":" + env.GetString("DB_PORT", "5432") + "/" + env.GetString("DB_NAME", "kicker_db") + "?sslmode=disable",
		},
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// DATABASE CONFIGURATION
	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	logger.Info("Connected to database %s", cfg.db.dsn)

	api := application{
		config: cfg,
		db: conn,
	}

	if err := api.run(api.mount()); err != nil {
		slog.Error("Serve has failed to start", "Error", err)
		os.Exit(1)
	}
}