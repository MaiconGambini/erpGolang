package suppliers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaiconGambini/erpGolang/backend/internal/suppliers"
)

func TestCreateInvalidJSON(t *testing.T) {
	handler := suppliers.NewHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/suppliers", bytes.NewBufferString("{"))
	rec := httptest.NewRecorder()
	handler.Create(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
