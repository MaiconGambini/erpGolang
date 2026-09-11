package auth_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaiconGambini/erpGolang/backend/internal/auth"
)

func TestLoginInvalidJSON(t *testing.T) {
	handler := auth.NewHandler(nil, 12, false)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString("{"))
	rec := httptest.NewRecorder()
	handler.Login(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
