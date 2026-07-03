package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MaiconGambini/erpGolang/backend/internal/app"
	"github.com/MaiconGambini/erpGolang/backend/internal/audit"
	"github.com/MaiconGambini/erpGolang/backend/internal/auth"
	"github.com/MaiconGambini/erpGolang/backend/internal/config"
	"github.com/MaiconGambini/erpGolang/backend/internal/products"
	"github.com/MaiconGambini/erpGolang/backend/internal/platform/logger"
	redisplatform "github.com/MaiconGambini/erpGolang/backend/internal/platform/redis"
	"github.com/MaiconGambini/erpGolang/backend/internal/platform/validation"
	"github.com/MaiconGambini/erpGolang/backend/internal/tenants"
	"github.com/MaiconGambini/erpGolang/backend/internal/users"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.AppEnv)

	ctx := context.Background()
	db, err := database.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Error("database connection failed", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()

	redisClient, err := redisplatform.Open(ctx, cfg.RedisURL)
	if err != nil {
		log.Error("redis connection failed", slog.Any("error", err))
		os.Exit(1)
	}
	defer func() { _ = redisClient.Close() }()

	validator := validation.New()
	auditRecorder := audit.NewService(db)

	deps := app.Dependencies{
		Config:    cfg,
		DB:        db,
		Redis:     redisClient,
		Validator: validator,
		Audit:     auditRecorder,
		Logger:    log,
		Modules: []app.Module{
			auth.NewModule(),
			users.NewModule(),
			tenants.NewModule(),
			customers.NewModule(),
			products.NewModule(),
			audit.NewModule(),
		},
	}

	handler := app.NewRouter(deps)
	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Info("starting goERP API", slog.String("addr", cfg.HTTPAddr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server stopped unexpectedly", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	shutdownCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-shutdownCtx.Done()

	shutdown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdown); err != nil {
		log.Error("graceful shutdown failed", slog.Any("error", err))
		os.Exit(1)
	}
}
