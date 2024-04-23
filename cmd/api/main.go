package main

import (
	bookshelf "bookshelf-api"
	"bookshelf-api/pkg/config"
	"bookshelf-api/pkg/handler"
	"bookshelf-api/pkg/service"
	"bookshelf-api/pkg/storage"
	"bookshelf-api/pkg/storage/postgres"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Info(
		"starting bookshelf",
		slog.String("env", cfg.Env),
		slog.String("version", "1.0"),
	)
	db, err := postgres.New(cfg)
	if err != nil {
		log.Error("failed to init storage", slog.String("err", err.Error()))
		os.Exit(1)
	}
	stor := storage.New(db)
	services := service.New(stor)
	handlers := handler.New(services)

	log.Info("starting server", slog.String("address", cfg.Address))
	srv := new(bookshelf.Server)
	if err := srv.Run(cfg, handlers.InitRoutes(log)); err != nil {
		log.Error("failed to start server")
	}
	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
