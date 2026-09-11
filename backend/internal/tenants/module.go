package tenants

import (
	"net/http"

	"github.com/MaiconGambini/erpGolang/backend/gen/db"
	"github.com/MaiconGambini/erpGolang/backend/internal/app"
	"github.com/MaiconGambini/erpGolang/backend/internal/auth"
	"github.com/MaiconGambini/erpGolang/backend/internal/platform/middleware"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/httpx"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/pgutil"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/tenantctx"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Module struct{}

func NewModule() Module { return Module{} }

func (Module) Name() string { return "tenants" }

func (Module) Register(r chi.Router, deps app.Deps) {
	queries := db.New(deps.DB)
	authSvc := auth.NewService(deps.DB, deps.Config)

	r.Route("/tenants", func(tenants chi.Router) {
		tenants.Use(middleware.AuthJWT(authSvc.JWT()))
		tenants.Use(middleware.TenantScope)
		tenants.Get("/current", func(w http.ResponseWriter, req *http.Request) {
			tenantID, _ := tenantctx.TenantIDFromContext(req.Context())
			tid, _ := uuid.Parse(tenantID)
			tenant, err := queries.GetTenantByID(req.Context(), pgutil.UUIDToPg(tid))
			if err != nil {
				httpx.Error(w, "NOT_FOUND", "tenant not found", http.StatusNotFound)
				return
			}
			id, _ := pgutil.PgToUUID(tenant.ID)
			httpx.JSON(w, http.StatusOK, map[string]any{
				"id":     id.String(),
				"slug":   tenant.Slug,
				"name":   tenant.Name,
				"status": tenant.Status,
			})
		})
	})
}
