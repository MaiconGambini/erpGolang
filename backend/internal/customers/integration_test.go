//go:build integration

package customers_test

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
	"github.com/MaiconGambini/erpGolang/backend/internal/platform/database"
	"github.com/MaiconGambini/erpGolang/backend/internal/platform/validation"
	redisplatform "github.com/MaiconGambini/erpGolang/backend/internal/platform/redis"
	"github.com/MaiconGambini/erpGolang/backend/internal/tenants"
	"github.com/MaiconGambini/erpGolang/backend/internal/users"
)

// Integration test: requires DATABASE_URL, JWT secrets, and seeded tenants.
// Run: go test -tags=integration ./internal/customers/... -count=1

func TestCrossTenantCustomerReturnsNotFound(t *testing.T) {
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
			audit.NewModule(),
		},
	}

	handler := app.NewRouter(deps)
	server := httptest.NewServer(handler)
	defer server.Close()

	acmeToken := login(t, server.URL, "acme", "admin@acme.com", "admin123")
	betaToken := login(t, server.URL, "beta", "admin@beta.com", "admin123")

	body := `{"name":"Isolation Test Customer","active":true}`
	req, _ := http.NewRequest(http.MethodPost, server.URL+"/api/v1/customers", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+acmeToken)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d", res.StatusCode)
	}

	var created struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}

	getReq, _ := http.NewRequest(http.MethodGet, server.URL+"/api/v1/customers/"+created.Data.ID, nil)
	getReq.Header.Set("Authorization", "Bearer "+betaToken)
	getRes, err := http.DefaultClient.Do(getReq)
	if err != nil {
		t.Fatal(err)
	}
	defer getRes.Body.Close()
	if getRes.StatusCode != http.StatusNotFound {
		t.Fatalf("cross-tenant get: expected 404, got %d", getRes.StatusCode)
	}
}

func login(t *testing.T, baseURL, slug, email, password string) string {
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
