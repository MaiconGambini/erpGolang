package reports

import (
	"github.com/MaiconGambini/erpGolang/backend/internal/app"
	"github.com/MaiconGambini/erpGolang/backend/internal/auth"
	"github.com/MaiconGambini/erpGolang/backend/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
)

type Module struct{}

func NewModule() Module { return Module{} }

func (Module) Name() string { return "reports" }

func (Module) Register(r chi.Router, deps app.Deps) {
	svc := NewService(deps.DB)
	handler := NewHandler(svc)
	authSvc := auth.NewService(deps.DB, deps.Config)
	financial := middleware.RequireRole("admin", "manager")
	read := middleware.RequireRole("admin", "manager", "operator", "viewer")

	r.Route("/reports", func(reports chi.Router) {
		reports.Use(middleware.AuthJWT(authSvc.JWT()))
		reports.Use(middleware.TenantScope)
		reports.With(financial).Get("/sales-by-day", handler.SalesByDay)
		reports.With(financial).Get("/top-products", handler.TopProducts)
		reports.With(financial).Get("/sales-summary.pdf", handler.SalesSummaryPDF)
		reports.With(read).Get("/sales/{id}/pdf", handler.SalePDF)
	})
}
