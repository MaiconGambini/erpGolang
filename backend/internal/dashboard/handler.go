package dashboard

import (
	"net/http"
	"strconv"

	"github.com/MaiconGambini/erpGolang/backend/internal/shared/authctx"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/httpx"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/inventory"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/tenantctx"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Summary(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := tenantctx.TenantIDFromContext(r.Context())
	thresholdRaw, _ := strconv.Atoi(r.URL.Query().Get("threshold"))
	threshold := inventory.NormalizeLowStockThreshold(thresholdRaw)
	summary, err := h.svc.Summary(r.Context(), SummaryParams{
		TenantID:  tenantID,
		Threshold: threshold,
	})
	if err != nil {
		httpx.Error(w, "INTERNAL_ERROR", "failed to load dashboard summary", http.StatusInternalServerError)
		return
	}
	if user, err := authctx.UserFromContext(r.Context()); err == nil {
		if user.Role != "admin" && user.Role != "manager" {
			summary.ConfirmedSalesCount = 0
			summary.ConfirmedSalesTotal = "0"
		}
	}
	httpx.JSON(w, http.StatusOK, summary)
}
