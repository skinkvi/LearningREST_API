package main

import (
	"log/slog"
	"net/http"
	"os"
	"rest_api_app/internal/config"
	"rest_api_app/internal/handlers/delete"
	"rest_api_app/internal/handlers/redirect"
	"rest_api_app/internal/handlers/url/save"
	"rest_api_app/internal/lib/logger/sl"
	"rest_api_app/internal/storage/postgres"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

const (
	envLocal = "local"
	envDev   = "dev"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("starting url-shortner", slog.String("env", cfg.Env))

	storage, err := postgres.New()
	if err != nil {
		log.Error("falied to init storage", sl.Err(err))
		os.Exit(1)
	}

	if err := storage.AutoMigrationTablePg(); err != nil {
		log.Error("falied to migrate database", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(log, storage))
	router.Get("/{alias}", redirect.New(log, storage))
	router.Post("/delete?alias={alias}", delete.New(log, storage))

	log.Info("server starting", slog.String("on ", cfg.HTTPServer.Address))

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("falied to start server")
	}

	log.Error("server stopped")

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return log
}
