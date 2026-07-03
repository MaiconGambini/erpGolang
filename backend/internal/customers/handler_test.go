package customers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaiconGambini/erpGolang/backend/internal/customers"
)

func TestCreateInvalidJSON(t *testing.T) {
	handler := customers.NewHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewBufferString("{"))
	rec := httptest.NewRecorder()
	handler.Create(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
