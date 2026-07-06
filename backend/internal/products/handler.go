package products

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/MaiconGambini/erpGolang/backend/internal/shared/authctx"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/export"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/httpx"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/inventory"
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
	if q.Get("format") == "csv" {
		items, _, err := h.svc.List(r.Context(), ListParams{
			TenantID: tenantID, Search: q.Get("search"), Active: ParseBoolQuery(q.Get("active")),
			Limit: 10000, Offset: 0,
		})
		if err != nil {
			httpx.Error(w, "INTERNAL_ERROR", "failed to export products", http.StatusInternalServerError)
			return
		}
		rows := make([][]string, 0, len(items))
		for _, p := range items {
			active := "false"
			if p.Active {
				active = "true"
			}
			rows = append(rows, []string{p.Name, p.Sku, p.Price, strconv.Itoa(p.Stock), active})
		}
		_ = export.WriteCSV(w, "produtos.csv", []string{"nome", "sku", "preco", "estoque", "ativo"}, rows)
		return
	}
	items, total, err := h.svc.List(r.Context(), ListParams{
		TenantID: tenantID,
		Search:   q.Get("search"),
		Active:   ParseBoolQuery(q.Get("active")),
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		httpx.Error(w, "INTERNAL_ERROR", "failed to list products", http.StatusInternalServerError)
		return
	}
	httpx.Paginated(w, items, httpx.PaginationMeta{
		Limit: limit, Offset: offset, Total: total,
	})
}

func (h *Handler) LowStock(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := tenantctx.TenantIDFromContext(r.Context())
	q := r.URL.Query()
	thresholdRaw, _ := strconv.Atoi(q.Get("threshold"))
	threshold := int(inventory.NormalizeLowStockThreshold(thresholdRaw))
	limit, _ := strconv.Atoi(q.Get("limit"))
	items, err := h.svc.LowStock(r.Context(), LowStockParams{
		TenantID:  tenantID,
		Threshold: threshold,
		Limit:     limit,
	})
	if err != nil {
		httpx.Error(w, "INTERNAL_ERROR", "failed to list low stock products", http.StatusInternalServerError)
		return
	}
	httpx.JSON(w, http.StatusOK, items)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := tenantctx.TenantIDFromContext(r.Context())
	id := chi.URLParam(r, "id")
	item, err := h.svc.Get(r.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, errNotFound) {
			httpx.Error(w, "NOT_FOUND", "product not found", http.StatusNotFound)
			return
		}
		httpx.Error(w, "INTERNAL_ERROR", "failed to get product", http.StatusInternalServerError)
		return
	}
	httpx.JSON(w, http.StatusOK, item)
}

type productRequest struct {
	Name    string  `json:"name"`
	Sku     string  `json:"sku"`
	Price   string  `json:"price"`
	Stock   int     `json:"stock"`
	Unit    string  `json:"unit"`
	Barcode *string `json:"barcode"`
	Active  bool    `json:"active"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := tenantctx.TenantIDFromContext(r.Context())
	user, _ := authctx.UserFromContext(r.Context())
	var req productRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, "INVALID_JSON", "invalid request body", http.StatusBadRequest)
		return
	}
	item, err := h.svc.Create(r.Context(), tenantID, user.ID, CreateInput{
		Name: req.Name, Sku: req.Sku, Price: req.Price, Stock: req.Stock,
		Unit: req.Unit, Barcode: req.Barcode, Active: req.Active,
	})
	if err != nil {
		writeServiceError(w, err, "failed to create product")
		return
	}
	httpx.JSON(w, http.StatusCreated, item)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := tenantctx.TenantIDFromContext(r.Context())
	user, _ := authctx.UserFromContext(r.Context())
	id := chi.URLParam(r, "id")
	var req productRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, "INVALID_JSON", "invalid request body", http.StatusBadRequest)
		return
	}
	item, err := h.svc.Update(r.Context(), tenantID, user.ID, id, CreateInput{
		Name: req.Name, Sku: req.Sku, Price: req.Price, Stock: req.Stock,
		Unit: req.Unit, Barcode: req.Barcode, Active: req.Active,
	})
	if err != nil {
		writeServiceError(w, err, "failed to update product")
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
			httpx.Error(w, "NOT_FOUND", "product not found", http.StatusNotFound)
			return
		}
		httpx.Error(w, "INTERNAL_ERROR", "failed to delete product", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeServiceError(w http.ResponseWriter, err error, fallback string) {
	switch {
	case errors.Is(err, errNotFound):
		httpx.Error(w, "NOT_FOUND", "product not found", http.StatusNotFound)
	case errors.Is(err, errValidation):
		httpx.Error(w, "VALIDATION_ERROR", err.Error(), http.StatusBadRequest)
	case errors.Is(err, errDuplicateSKU):
		httpx.Error(w, "DUPLICATE_SKU", "sku already exists for this tenant", http.StatusConflict)
	default:
		httpx.Error(w, "INTERNAL_ERROR", fallback, http.StatusInternalServerError)
	}
}
