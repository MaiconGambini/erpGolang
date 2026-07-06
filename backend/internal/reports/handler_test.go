//go:build integration

package reports_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/MaiconGambini/erpGolang/backend/internal/app"
	"github.com/MaiconGambini/erpGolang/backend/internal/audit"
	"github.com/MaiconGambini/erpGolang/backend/internal/auth"
	"github.com/MaiconGambini/erpGolang/backend/internal/config"
	"github.com/MaiconGambini/erpGolang/backend/internal/customers"
	"github.com/MaiconGambini/erpGolang/backend/internal/platform/database"
	"github.com/MaiconGambini/erpGolang/backend/internal/platform/logger"
	"github.com/MaiconGambini/erpGolang/backend/internal/platform/validation"
	redisplatform "github.com/MaiconGambini/erpGolang/backend/internal/platform/redis"
	"github.com/MaiconGambini/erpGolang/backend/internal/products"
	"github.com/MaiconGambini/erpGolang/backend/internal/reports"
	"github.com/MaiconGambini/erpGolang/backend/internal/sales"
	"github.com/MaiconGambini/erpGolang/backend/internal/tenants"
	"github.com/MaiconGambini/erpGolang/backend/internal/users"
)

func TestCustomersCSVExport(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	server := newTestServer(t)
	defer server.Close()
	token := login(t, server.URL)

	req, _ := http.NewRequest(http.MethodGet, server.URL+"/api/v1/customers?format=csv", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("csv export: expected 200, got %d", res.StatusCode)
	}
	ct := res.Header.Get("Content-Type")
	if !strings.Contains(ct, "text/csv") {
		t.Fatalf("expected text/csv content-type, got %q", ct)
	}
}

func TestSalePDFReturnsPDF(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	server := newTestServer(t)
	defer server.Close()
	token := login(t, server.URL)

	customerID := createCustomer(t, server.URL, token)
	sku := fmt.Sprintf("PDF-%d", time.Now().UnixNano())
	productID := createProduct(t, server.URL, token, sku)
	saleID := createAndConfirmSale(t, server.URL, token, customerID, productID)

	req, _ := http.NewRequest(http.MethodGet, server.URL+"/api/v1/reports/sales/"+saleID+"/pdf", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("sale pdf: expected 200, got %d", res.StatusCode)
	}
	if ct := res.Header.Get("Content-Type"); ct != "application/pdf" {
		t.Fatalf("expected application/pdf, got %q", ct)
	}
}

func TestSalesSummaryPDFReturnsPDF(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	server := newTestServer(t)
	defer server.Close()
	token := login(t, server.URL)

	req, _ := http.NewRequest(http.MethodGet, server.URL+"/api/v1/reports/sales-summary.pdf", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("summary pdf: expected 200, got %d", res.StatusCode)
	}
	if ct := res.Header.Get("Content-Type"); ct != "application/pdf" {
		t.Fatalf("expected application/pdf, got %q", ct)
	}
}

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	cfg := config.Load()
	ctx := context.Background()
	db, err := database.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		t.Skipf("database unavailable: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	redisClient, err := redisplatform.Open(ctx, cfg.RedisURL)
	if err != nil {
		t.Skipf("redis unavailable: %v", err)
	}
	t.Cleanup(func() { _ = redisClient.Close() })

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
			products.NewModule(),
			sales.NewModule(),
			reports.NewModule(),
			audit.NewModule(),
		},
	}
	return httptest.NewServer(app.NewRouter(deps))
}

func login(t *testing.T, baseURL string) string {
	t.Helper()
	payload := map[string]string{"tenantSlug": "acme", "email": "admin@acme.com", "password": "admin123"}
	b, _ := json.Marshal(payload)
	res, err := http.Post(baseURL+"/api/v1/auth/login", "application/json", bytes.NewBuffer(b))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
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

func createCustomer(t *testing.T, baseURL, token string) string {
	t.Helper()
	body := fmt.Sprintf(`{"name":"PDF Customer %d","active":true}`, time.Now().UnixNano())
	req, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/customers", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var out struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	_ = json.NewDecoder(res.Body).Decode(&out)
	return out.Data.ID
}

func createProduct(t *testing.T, baseURL, token, sku string) string {
	t.Helper()
	body := fmt.Sprintf(`{"name":"PDF Product","sku":%q,"price":"15.00","stock":5,"unit":"un","active":true}`, sku)
	req, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/products", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var out struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	_ = json.NewDecoder(res.Body).Decode(&out)
	return out.Data.ID
}

func createAndConfirmSale(t *testing.T, baseURL, token, customerID, productID string) string {
	t.Helper()
	body := fmt.Sprintf(`{"customerId":%q,"items":[{"productId":%q,"quantity":1}]}`, customerID, productID)
	req, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/sales", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var out struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	_ = json.NewDecoder(res.Body).Decode(&out)

	conf, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/sales/"+out.Data.ID+"/confirm", nil)
	conf.Header.Set("Authorization", "Bearer "+token)
	confRes, err := http.DefaultClient.Do(conf)
	if err != nil {
		t.Fatal(err)
	}
	_ = confRes.Body.Close()
	return out.Data.ID
}
