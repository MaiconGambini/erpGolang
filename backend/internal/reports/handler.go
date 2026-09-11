package reports

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/MaiconGambini/erpGolang/backend/internal/shared/httpx"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/tenantctx"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) SalesByDay(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := tenantctx.TenantIDFromContext(r.Context())
	from, to := ParseRangeQuery(r.URL.Query().Get("from"), r.URL.Query().Get("to"))
	items, err := h.svc.SalesByDay(r.Context(), RangeParams{TenantID: tenantID, From: from, To: to})
	if err != nil {
		httpx.Error(w, "INTERNAL_ERROR", "failed to load sales by day", http.StatusInternalServerError)
		return
	}
	httpx.JSON(w, http.StatusOK, items)
}

func (h *Handler) TopProducts(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := tenantctx.TenantIDFromContext(r.Context())
	from, to := ParseRangeQuery(r.URL.Query().Get("from"), r.URL.Query().Get("to"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	items, err := h.svc.TopProducts(r.Context(), RangeParams{
		TenantID: tenantID, From: from, To: to, Limit: int32(limit),
	})
	if err != nil {
		httpx.Error(w, "INTERNAL_ERROR", "failed to load top products", http.StatusInternalServerError)
		return
	}
	httpx.JSON(w, http.StatusOK, items)
}

func (h *Handler) SalePDF(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := tenantctx.TenantIDFromContext(r.Context())
	id := chi.URLParam(r, "id")
	data, err := h.svc.SalePDF(r.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.Error(w, "NOT_FOUND", "sale not found", http.StatusNotFound)
			return
		}
		httpx.Error(w, "INTERNAL_ERROR", "failed to generate pdf", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="venda-`+id+`.pdf"`)
	_, _ = w.Write(data)
}

func (h *Handler) SalesSummaryPDF(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := tenantctx.TenantIDFromContext(r.Context())
	from, to := ParseRangeQuery(r.URL.Query().Get("from"), r.URL.Query().Get("to"))
	data, err := h.svc.SalesSummaryPDF(r.Context(), tenantID, from, to)
	if err != nil {
		httpx.Error(w, "INTERNAL_ERROR", "failed to generate pdf", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="vendas-resumo.pdf"`)
	_, _ = w.Write(data)
}
