package audit

import (
	"net/http"
	"strconv"

	"github.com/MaiconGambini/erpGolang/backend/internal/shared/httpx"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/tenantctx"
)

type ListHandler struct {
	svc *ListService
}

func NewListHandler(svc *ListService) *ListHandler {
	return &ListHandler{svc: svc}
}

func (h *ListHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := tenantctx.TenantIDFromContext(r.Context())
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit <= 0 {
		limit = 50
	}
	offset, _ := strconv.Atoi(q.Get("offset"))
	items, total, err := h.svc.List(r.Context(), ListParams{
		TenantID: tenantID, Limit: limit, Offset: offset,
	})
	if err != nil {
		httpx.Error(w, "INTERNAL_ERROR", "failed to list audit logs", http.StatusInternalServerError)
		return
	}
	httpx.Paginated(w, items, httpx.PaginationMeta{Limit: limit, Offset: offset, Total: total})
}
