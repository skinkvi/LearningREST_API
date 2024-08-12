package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"rest_api_app/internal/client/sso/grpc"
	"rest_api_app/internal/config"
	"rest_api_app/internal/handlers/delete"
	"rest_api_app/internal/handlers/redirect"
	"rest_api_app/internal/handlers/url/save"
	"rest_api_app/internal/lib/logger/sl"
	"rest_api_app/internal/storage/postgres"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	ssov1 "github.com/skinkvi/protosSTT/gen/go/sso"
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

	// Инициализация gRPC клиента для SSO
	ssoClient, err := grpc.New(context.Background(), log, cfg.SSO.Address, cfg.SSO.Timeout, cfg.SSO.RetriesCount)
	if err != nil {
		log.Error("failed to initialize SSO gRPC client", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// Использование SSO для аутентификации и авторизации
	router.Route("/url", func(r chi.Router) {
		r.Use(ssoMiddleware(ssoClient, log))
		r.Post("/", save.New(log, storage))
		r.Delete("/url/{alias}", delete.New(log, storage))
	})

	router.Get("/{alias}", redirect.New(log, storage))

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

func ssoMiddleware(ssoClient *grpc.Client, log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
				return
			}

			// Вызов SSO для проверки токена
			resp, err := ssoClient.Api.ValidateToken(context.Background(), &ssov1.ValidateTokenRequest{Token: token})
			if err != nil {
				log.Error("failed to validate token", sl.Err(err))
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			if !resp.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// TODO: переопределить с типа стринг на собственный тип что бы избежать колизии именно поэтому он и подчеркивает и готовит Warning
			ctx := context.WithValue(r.Context(), "user", resp.User)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
