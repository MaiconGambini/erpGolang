package products_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaiconGambini/erpGolang/backend/internal/products"
)

func TestCreateInvalidJSON(t *testing.T) {
	handler := products.NewHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString("{"))
	rec := httptest.NewRecorder()
	handler.Create(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
