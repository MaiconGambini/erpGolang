package users

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/MaiconGambini/erpGolang/backend/internal/shared/authctx"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/httpx"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/tenantctx"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := tenantctx.TenantIDFromContext(r.Context())
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit <= 0 {
		limit = 20
	}
	offset, _ := strconv.Atoi(q.Get("offset"))
	items, total, err := h.svc.List(r.Context(), ListParams{
		TenantID: tenantID, Limit: limit, Offset: offset,
	})
	if err != nil {
		httpx.Error(w, "INTERNAL_ERROR", "failed to list users", http.StatusInternalServerError)
		return
	}
	httpx.Paginated(w, items, httpx.PaginationMeta{Limit: limit, Offset: offset, Total: total})
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := tenantctx.TenantIDFromContext(r.Context())
	id := chi.URLParam(r, "id")
	item, err := h.svc.Get(r.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, errNotFound) {
			httpx.Error(w, "NOT_FOUND", "user not found", http.StatusNotFound)
			return
		}
		httpx.Error(w, "INTERNAL_ERROR", "failed to get user", http.StatusInternalServerError)
		return
	}
	httpx.JSON(w, http.StatusOK, item)
}

type userRequest struct {
	Name   string `json:"name"`
	Role   string `json:"role"`
	Active bool   `json:"active"`
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := tenantctx.TenantIDFromContext(r.Context())
	actor, _ := authctx.UserFromContext(r.Context())
	id := chi.URLParam(r, "id")
	var req userRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, "INVALID_JSON", "invalid request body", http.StatusBadRequest)
		return
	}
	item, err := h.svc.Update(r.Context(), tenantID, actor.ID, id, UpdateInput(req))
	if err != nil {
		switch {
		case errors.Is(err, errNotFound):
			httpx.Error(w, "NOT_FOUND", "user not found", http.StatusNotFound)
		case errors.Is(err, errValidation):
			httpx.Error(w, "VALIDATION_ERROR", "invalid user data", http.StatusBadRequest)
		case errors.Is(err, errSelfDisable):
			httpx.Error(w, "VALIDATION_ERROR", "cannot deactivate your own account", http.StatusBadRequest)
		default:
			httpx.Error(w, "INTERNAL_ERROR", "failed to update user", http.StatusInternalServerError)
		}
		return
	}
	httpx.JSON(w, http.StatusOK, item)
}
