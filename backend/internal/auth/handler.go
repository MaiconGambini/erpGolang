package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/MaiconGambini/erpGolang/backend/internal/shared/authctx"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/httpx"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/tenantctx"
)

const refreshCookieName = "refresh_token"

type Handler struct {
	svc           *Service
	bcryptCost    int
	secureCookies bool
}

func NewHandler(svc *Service, bcryptCost int, secureCookies bool) *Handler {
	return &Handler{svc: svc, bcryptCost: bcryptCost, secureCookies: secureCookies}
}

type loginRequest struct {
	TenantSlug string `json:"tenantSlug"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, "INVALID_JSON", "invalid request body", http.StatusBadRequest)
		return
	}
	result, err := h.svc.Login(r.Context(), LoginInput{
		TenantSlug: req.TenantSlug,
		Email:      req.Email,
		Password:   req.Password,
	}, h.bcryptCost)
	if err != nil {
		httpx.Error(w, "UNAUTHORIZED", "invalid credentials", http.StatusUnauthorized)
		return
	}
	h.setRefreshCookie(w, result.RefreshToken)
	httpx.JSON(w, http.StatusOK, map[string]any{
		"accessToken": result.AccessToken,
		"user": map[string]string{
			"id":       result.UserID,
			"tenantId": result.TenantID,
			"email":    result.Email,
			"name":     result.Name,
			"role":     result.Role,
		},
	})
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie(refreshCookieName)
	if err != nil || token.Value == "" {
		httpx.Error(w, "UNAUTHORIZED", "refresh token missing", http.StatusUnauthorized)
		return
	}
	result, err := h.svc.Refresh(r.Context(), token.Value)
	if err != nil {
		h.clearRefreshCookie(w)
		httpx.Error(w, "UNAUTHORIZED", "invalid refresh token", http.StatusUnauthorized)
		return
	}
	h.setRefreshCookie(w, result.RefreshToken)
	httpx.JSON(w, http.StatusOK, map[string]any{
		"accessToken": result.AccessToken,
		"user": map[string]string{
			"id":       result.UserID,
			"tenantId": result.TenantID,
			"email":    result.Email,
			"name":     result.Name,
			"role":     result.Role,
		},
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	token, _ := r.Cookie(refreshCookieName)
	if token != nil {
		_ = h.svc.Logout(r.Context(), token.Value)
	}
	h.clearRefreshCookie(w)
	httpx.JSON(w, http.StatusOK, map[string]string{"status": "logged_out"})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	user, err := authctx.UserFromContext(r.Context())
	if err != nil {
		httpx.Error(w, "UNAUTHORIZED", "unauthorized", http.StatusUnauthorized)
		return
	}
	tenantID, err := tenantctx.TenantIDFromContext(r.Context())
	if err != nil {
		httpx.Error(w, "UNAUTHORIZED", "unauthorized", http.StatusUnauthorized)
		return
	}
	profile, err := h.svc.Me(r.Context(), user.ID, tenantID)
	if err != nil {
		httpx.Error(w, "NOT_FOUND", "user not found", http.StatusNotFound)
		return
	}
	httpx.JSON(w, http.StatusOK, profile)
}

func (h *Handler) setRefreshCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshCookieName,
		Value:    token,
		Path:     "/api/v1/auth",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   h.secureCookies,
		MaxAge:   int((30 * 24 * time.Hour).Seconds()),
	})
}

func (h *Handler) clearRefreshCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshCookieName,
		Value:    "",
		Path:     "/api/v1/auth",
		HttpOnly: true,
		Secure:   h.secureCookies,
		MaxAge:   -1,
	})
}

