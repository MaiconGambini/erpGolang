package sales_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaiconGambini/erpGolang/backend/internal/sales"
)

func TestCreateInvalidJSON(t *testing.T) {
	handler := sales.NewHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/sales", bytes.NewBufferString("{"))
	rec := httptest.NewRecorder()
	handler.Create(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
