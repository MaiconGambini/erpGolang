//go:build integration

package dashboard_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaiconGambini/erpGolang/backend/internal/app"
	"github.com/MaiconGambini/erpGolang/backend/internal/audit"
	"github.com/MaiconGambini/erpGolang/backend/internal/auth"
	"github.com/MaiconGambini/erpGolang/backend/internal/config"
	"github.com/MaiconGambini/erpGolang/backend/internal/customers"
	"github.com/MaiconGambini/erpGolang/backend/internal/dashboard"
	"github.com/MaiconGambini/erpGolang/backend/internal/platform/database"
	"github.com/MaiconGambini/erpGolang/backend/internal/platform/validation"
	redisplatform "github.com/MaiconGambini/erpGolang/backend/internal/platform/redis"
	"github.com/MaiconGambini/erpGolang/backend/internal/tenants"
	"github.com/MaiconGambini/erpGolang/backend/internal/users"
)

func TestDashboardSummaryTenantScoped(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	cfg := config.Load()
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
	defer redisClient.Close()

	deps := app.Dependencies{
		Config:    cfg,
		DB:        db,
		Redis:     redisClient,
		Validator: validation.New(),
		Audit:     audit.NewService(db),
		Modules: []app.Module{
			auth.NewModule(),
			users.NewModule(),
			tenants.NewModule(),
			customers.NewModule(),
			dashboard.NewModule(),
			audit.NewModule(),
		},
	}

	server := httptest.NewServer(app.NewRouter(deps))
	defer server.Close()

	acmeToken := loginIntegration(t, server.URL, "acme", "admin@acme.com", "admin123")
	betaToken := loginIntegration(t, server.URL, "beta", "admin@beta.com", "admin123")

	acmeBefore := fetchSummary(t, server.URL, acmeToken)
	betaBefore := fetchSummary(t, server.URL, betaToken)

	body := `{"name":"Dashboard Isolation Customer","active":true}`
	req, _ := http.NewRequest(http.MethodPost, server.URL+"/api/v1/customers", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+acmeToken)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("create customer: expected 201, got %d", res.StatusCode)
	}

	acmeAfter := fetchSummary(t, server.URL, acmeToken)
	betaAfter := fetchSummary(t, server.URL, betaToken)

	if acmeAfter.ActiveCustomers != acmeBefore.ActiveCustomers+1 {
		t.Fatalf("acme active_customers: before=%d after=%d", acmeBefore.ActiveCustomers, acmeAfter.ActiveCustomers)
	}
	if betaAfter.ActiveCustomers != betaBefore.ActiveCustomers {
		t.Fatalf("beta active_customers changed: before=%d after=%d", betaBefore.ActiveCustomers, betaAfter.ActiveCustomers)
	}
}

type summaryCounts struct {
	ActiveCustomers     int64 `json:"activeCustomers"`
	ConfirmedSalesCount int64 `json:"confirmedSalesCount"`
}

func fetchSummary(t *testing.T, baseURL, token string) summaryCounts {
	t.Helper()
	req, _ := http.NewRequest(http.MethodGet, baseURL+"/api/v1/dashboard/summary", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("summary: expected 200, got %d", res.StatusCode)
	}
	var out struct {
		Data summaryCounts `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	return out.Data
}

func loginIntegration(t *testing.T, baseURL, slug, email, password string) string {
	t.Helper()
	payload := map[string]string{"tenantSlug": slug, "email": email, "password": password}
	b, _ := json.Marshal(payload)
	res, err := http.Post(baseURL+"/api/v1/auth/login", "application/json", bytes.NewBuffer(b))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("login %s: expected 200, got %d", slug, res.StatusCode)
	}
	var out struct {
		Data struct {
			AccessToken string `json:"accessToken"`
		} `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	return out.Data.AccessToken
}
