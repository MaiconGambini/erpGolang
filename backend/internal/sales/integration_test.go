//go:build integration

package sales_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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
	"github.com/MaiconGambini/erpGolang/backend/internal/sales"
	"github.com/MaiconGambini/erpGolang/backend/internal/tenants"
	"github.com/MaiconGambini/erpGolang/backend/internal/users"
)

func TestConfirmDecrementsStock(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	server := newTestServer(t)
	defer server.Close()
	token := login(t, server.URL, "acme", "admin@acme.com", "admin123")

	customerID := createCustomer(t, server.URL, token, fmt.Sprintf("Sale Stock %d", time.Now().UnixNano()))
	sku := fmt.Sprintf("SKU-%d", time.Now().UnixNano())
	productID := createProduct(t, server.URL, token, sku, 10)
	saleID := createSale(t, server.URL, token, customerID, productID, 2)

	confirmSale(t, server.URL, token, saleID)

	stock := getProductStock(t, server.URL, token, productID)
	if stock != 8 {
		t.Fatalf("expected stock 8 after confirm, got %d", stock)
	}
}

func TestConfirmInsufficientStockReturns409(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	server := newTestServer(t)
	defer server.Close()
	token := login(t, server.URL, "acme", "admin@acme.com", "admin123")

	customerID := createCustomer(t, server.URL, token, fmt.Sprintf("Sale Low %d", time.Now().UnixNano()))
	sku := fmt.Sprintf("LOW-%d", time.Now().UnixNano())
	productID := createProduct(t, server.URL, token, sku, 1)
	saleID := createSale(t, server.URL, token, customerID, productID, 5)

	req, _ := http.NewRequest(http.MethodPost, server.URL+"/api/v1/sales/"+saleID+"/confirm", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409, got %d", res.StatusCode)
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
			audit.NewModule(),
		},
	}
	return httptest.NewServer(app.NewRouter(deps))
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

func createCustomer(t *testing.T, baseURL, token, name string) string {
	t.Helper()
	body := fmt.Sprintf(`{"name":%q,"active":true}`, name)
	req, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/customers", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("create customer: expected 201, got %d", res.StatusCode)
	}
	var out struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	return out.Data.ID
}

func createProduct(t *testing.T, baseURL, token, sku string, stock int) string {
	t.Helper()
	body := fmt.Sprintf(`{"name":"Product %s","sku":%q,"price":"10.00","stock":%d,"unit":"un","active":true}`, sku, sku, stock)
	req, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/products", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("create product: expected 201, got %d", res.StatusCode)
	}
	var out struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	return out.Data.ID
}

func createSale(t *testing.T, baseURL, token, customerID, productID string, qty int) string {
	t.Helper()
	body := fmt.Sprintf(`{"customerId":%q,"items":[{"productId":%q,"quantity":%d}]}`, customerID, productID, qty)
	req, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/sales", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("create sale: expected 201, got %d", res.StatusCode)
	}
	var out struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	return out.Data.ID
}

func confirmSale(t *testing.T, baseURL, token, saleID string) {
	t.Helper()
	req, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/sales/"+saleID+"/confirm", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("confirm sale: expected 200, got %d", res.StatusCode)
	}
}

func getProductStock(t *testing.T, baseURL, token, productID string) int {
	t.Helper()
	req, _ := http.NewRequest(http.MethodGet, baseURL+"/api/v1/products/"+productID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("get product: expected 200, got %d", res.StatusCode)
	}
	var out struct {
		Data struct {
			Stock int `json:"stock"`
		} `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	return out.Data.Stock
}
