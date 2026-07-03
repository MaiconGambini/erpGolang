package customers

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
		TenantID: tenantID,
		Search:   q.Get("search"),
		Active:   ParseBoolQuery(q.Get("active")),
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		httpx.Error(w, "INTERNAL_ERROR", "failed to list customers", http.StatusInternalServerError)
		return
	}
	httpx.Paginated(w, items, httpx.PaginationMeta{
		Limit: limit, Offset: offset, Total: total,
	})
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := tenantctx.TenantIDFromContext(r.Context())
	id := chi.URLParam(r, "id")
	item, err := h.svc.Get(r.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, errNotFound) {
			httpx.Error(w, "NOT_FOUND", "customer not found", http.StatusNotFound)
			return
		}
		httpx.Error(w, "INTERNAL_ERROR", "failed to get customer", http.StatusInternalServerError)
		return
	}
	httpx.JSON(w, http.StatusOK, item)
}

type customerRequest struct {
	Name     string  `json:"name"`
	Document *string `json:"document"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Active   bool    `json:"active"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := tenantctx.TenantIDFromContext(r.Context())
	user, _ := authctx.UserFromContext(r.Context())
	var req customerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, "INVALID_JSON", "invalid request body", http.StatusBadRequest)
		return
	}
	item, err := h.svc.Create(r.Context(), tenantID, user.ID, CreateInput{
		Name: req.Name, Document: req.Document, Email: req.Email, Phone: req.Phone, Active: req.Active,
	})
	if err != nil {
		httpx.Error(w, "INTERNAL_ERROR", "failed to create customer", http.StatusInternalServerError)
		return
	}
	httpx.JSON(w, http.StatusCreated, item)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := tenantctx.TenantIDFromContext(r.Context())
	user, _ := authctx.UserFromContext(r.Context())
	id := chi.URLParam(r, "id")
	var req customerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, "INVALID_JSON", "invalid request body", http.StatusBadRequest)
		return
	}
	item, err := h.svc.Update(r.Context(), tenantID, user.ID, id, CreateInput{
		Name: req.Name, Document: req.Document, Email: req.Email, Phone: req.Phone, Active: req.Active,
	})
	if err != nil {
		if errors.Is(err, errNotFound) {
			httpx.Error(w, "NOT_FOUND", "customer not found", http.StatusNotFound)
			return
		}
		httpx.Error(w, "INTERNAL_ERROR", "failed to update customer", http.StatusInternalServerError)
		return
	}
	httpx.JSON(w, http.StatusOK, item)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := tenantctx.TenantIDFromContext(r.Context())
	user, _ := authctx.UserFromContext(r.Context())
	id := chi.URLParam(r, "id")
	if err := h.svc.Delete(r.Context(), tenantID, user.ID, id); err != nil {
		if errors.Is(err, errNotFound) {
			httpx.Error(w, "NOT_FOUND", "customer not found", http.StatusNotFound)
			return
		}
		httpx.Error(w, "INTERNAL_ERROR", "failed to delete customer", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
