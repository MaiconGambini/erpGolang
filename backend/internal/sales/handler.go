package sales

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/MaiconGambini/erpGolang/backend/internal/shared/authctx"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/export"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/httpx"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/querytime"
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
	from, to := parseOptionalDates(q.Get("from"), q.Get("to"))
	if q.Get("format") == "csv" {
		items, _, err := h.svc.List(r.Context(), ListParams{
			TenantID: tenantID, Search: q.Get("search"), Status: ParseStatusQuery(q.Get("status")),
			FromDate: from, ToDate: to, Limit: 10000, Offset: 0,
		})
		if err != nil {
			httpx.Error(w, "INTERNAL_ERROR", "failed to export sales", http.StatusInternalServerError)
			return
		}
		rows := make([][]string, 0, len(items))
		for _, s := range items {
			rows = append(rows, []string{s.CustomerName, s.Status, s.Total, s.CreatedAt})
		}
		_ = export.WriteCSV(w, "vendas.csv", []string{"cliente", "status", "total", "criado_em"}, rows)
		return
	}
	items, total, err := h.svc.List(r.Context(), ListParams{
		TenantID: tenantID,
		Search:   q.Get("search"),
		Status:   ParseStatusQuery(q.Get("status")),
		FromDate: from,
		ToDate:   to,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		httpx.Error(w, "INTERNAL_ERROR", "failed to list sales", http.StatusInternalServerError)
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
			httpx.Error(w, "NOT_FOUND", "sale not found", http.StatusNotFound)
			return
		}
		httpx.Error(w, "INTERNAL_ERROR", "failed to get sale", http.StatusInternalServerError)
		return
	}
	httpx.JSON(w, http.StatusOK, item)
}

type saleRequest struct {
	CustomerID string `json:"customerId"`
	Notes      *string `json:"notes"`
	Items      []struct {
		ProductID string `json:"productId"`
		Quantity  int    `json:"quantity"`
	} `json:"items"`
}

func parseSaleRequest(req saleRequest) CreateInput {
	items := make([]ItemInput, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, ItemInput{ProductID: it.ProductID, Quantity: it.Quantity})
	}
	return CreateInput{CustomerID: req.CustomerID, Notes: req.Notes, Items: items}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := tenantctx.TenantIDFromContext(r.Context())
	user, _ := authctx.UserFromContext(r.Context())
	var req saleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, "INVALID_JSON", "invalid request body", http.StatusBadRequest)
		return
	}
	item, err := h.svc.Create(r.Context(), tenantID, user.ID, parseSaleRequest(req))
	if err != nil {
		writeServiceError(w, err, "failed to create sale")
		return
	}
	httpx.JSON(w, http.StatusCreated, item)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := tenantctx.TenantIDFromContext(r.Context())
	user, _ := authctx.UserFromContext(r.Context())
	id := chi.URLParam(r, "id")
	var req saleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, "INVALID_JSON", "invalid request body", http.StatusBadRequest)
		return
	}
	item, err := h.svc.Update(r.Context(), tenantID, user.ID, id, parseSaleRequest(req))
	if err != nil {
		writeServiceError(w, err, "failed to update sale")
		return
	}
	httpx.JSON(w, http.StatusOK, item)
}

func (h *Handler) Confirm(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := tenantctx.TenantIDFromContext(r.Context())
	user, _ := authctx.UserFromContext(r.Context())
	id := chi.URLParam(r, "id")
	item, err := h.svc.Confirm(r.Context(), tenantID, user.ID, id)
	if err != nil {
		writeServiceError(w, err, "failed to confirm sale")
		return
	}
	httpx.JSON(w, http.StatusOK, item)
}

func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := tenantctx.TenantIDFromContext(r.Context())
	user, _ := authctx.UserFromContext(r.Context())
	id := chi.URLParam(r, "id")
	item, err := h.svc.Cancel(r.Context(), tenantID, user.ID, id)
	if err != nil {
		writeServiceError(w, err, "failed to cancel sale")
		return
	}
	httpx.JSON(w, http.StatusOK, item)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := tenantctx.TenantIDFromContext(r.Context())
	user, _ := authctx.UserFromContext(r.Context())
	id := chi.URLParam(r, "id")
	if err := h.svc.Delete(r.Context(), tenantID, user.ID, id); err != nil {
		writeServiceError(w, err, "failed to delete sale")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeServiceError(w http.ResponseWriter, err error, fallback string) {
	switch {
	case errors.Is(err, errNotFound):
		httpx.Error(w, "NOT_FOUND", "sale not found", http.StatusNotFound)
	case errors.Is(err, errValidation):
		httpx.Error(w, "VALIDATION_ERROR", err.Error(), http.StatusBadRequest)
	case errors.Is(err, errInvalidStatus):
		httpx.Error(w, "INVALID_STATUS", "invalid sale status for this operation", http.StatusConflict)
	case errors.Is(err, errInsufficientStock):
		httpx.Error(w, "INSUFFICIENT_STOCK", "insufficient product stock", http.StatusConflict)
	default:
		httpx.Error(w, "INTERNAL_ERROR", fallback, http.StatusInternalServerError)
	}
}

func parseOptionalDates(fromStr, toStr string) (*time.Time, *time.Time) {
	var from, to *time.Time
	if fromStr != "" {
		if t := querytime.ParseOptional(fromStr); t.Valid {
			v := t.Time
			from = &v
		}
	}
	if toStr != "" {
		if t := querytime.ParseOptional(toStr); t.Valid {
			v := t.Time
			to = &v
		}
	}
	return from, to
}
