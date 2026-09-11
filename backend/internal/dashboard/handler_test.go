package dashboard_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/MaiconGambini/erpGolang/backend/internal/app"
	"github.com/MaiconGambini/erpGolang/backend/internal/audit"
	"github.com/MaiconGambini/erpGolang/backend/internal/auth"
	"github.com/MaiconGambini/erpGolang/backend/internal/config"
	"github.com/MaiconGambini/erpGolang/backend/internal/customers"
	"github.com/MaiconGambini/erpGolang/backend/internal/dashboard"
	"github.com/MaiconGambini/erpGolang/backend/internal/platform/database"
	"github.com/MaiconGambini/erpGolang/backend/internal/platform/logger"
	redisplatform "github.com/MaiconGambini/erpGolang/backend/internal/platform/redis"
	"github.com/MaiconGambini/erpGolang/backend/internal/platform/validation"
	"github.com/MaiconGambini/erpGolang/backend/internal/tenants"
	"github.com/MaiconGambini/erpGolang/backend/internal/users"
)

func requireIntegrationEnv(t *testing.T) config.Config {
	t.Helper()
	if os.Getenv("DATABASE_URL") == "" {
		t.Skip("DATABASE_URL not set")
	}
	return config.Load()
}

func TestSummaryUnauthorized(t *testing.T) {
	cfg := requireIntegrationEnv(t)
	ctx := context.Background()
	db, err := database.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		t.Skipf("database unavailable: %v", err)
	}
	defer db.Close()

	redisClient, err := redisplatform.Open(ctx, cfg.RedisURL)
	if err != nil {
		t.Skipf("redis unavailable: %v", err)
	}
	defer func() { _ = redisClient.Close() }()

	deps := app.Dependencies{
		Config:    cfg,
		DB:        db,
		Redis:     redisClient,
		Validator: validation.New(),
		Audit:     audit.NewService(db),
		Logger:    logger.New("test"),
		Modules: []app.Module{
			auth.NewModule(),
			users.NewModule(),
			tenants.NewModule(),
			customers.NewModule(),
			dashboard.NewModule(),
			audit.NewModule(),
		},
	}

	handler := app.NewRouter(deps)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard/summary", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d body=%s", rec.Code, rec.Body.String())
	}
}
