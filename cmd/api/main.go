package main

import (
	bookshelf "bookshelf-api"
	"bookshelf-api/pkg/config"
	"bookshelf-api/pkg/handler"
	"bookshelf-api/pkg/service"
	"bookshelf-api/pkg/storage"
	"bookshelf-api/pkg/storage/postgres"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// @title Bookshelf API
// @version 1.0
// @description API server for bookshelf

// @host localhost:8081
// @BasePath /

// @securityDefinition.apikey ApiKeyAuth
// @in header
// @name Authorization

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
	repos := storage.New(db)
	services := service.New(repos)
	handlers := handler.New(services)

	log.Info("starting server", slog.String("address", cfg.Address))
	srv := new(bookshelf.Server)
	go func() {
		if err := srv.Run(cfg, handlers.InitRoutes(log)); err != nil {
			log.Error("failed to start server")
		}
	}()
	log.Info("server started", slog.String("address", cfg.Address))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info("shutting down server")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Error("server shutdown failed", slog.String("err", err.Error()))
	}

	if err := db.Close(); err != nil {
		log.Error("failed to close database", slog.String("err", err.Error()))
	}
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
