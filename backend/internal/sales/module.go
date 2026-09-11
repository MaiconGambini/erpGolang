package sales

import (
	"github.com/MaiconGambini/erpGolang/backend/internal/app"
	"github.com/MaiconGambini/erpGolang/backend/internal/auth"
	"github.com/MaiconGambini/erpGolang/backend/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
)

type Module struct{}

func NewModule() Module { return Module{} }

func (Module) Name() string { return "sales" }

func (Module) Register(r chi.Router, deps app.Deps) {
	svc := NewService(deps.DB, deps.Audit)
	handler := NewHandler(svc)
	authSvc := auth.NewService(deps.DB, deps.Config)
	read := middleware.RequireRole("admin", "manager", "operator", "viewer")
	write := middleware.RequireRole("admin", "manager", "operator")

	r.Route("/sales", func(sales chi.Router) {
		sales.Use(middleware.AuthJWT(authSvc.JWT()))
		sales.Use(middleware.TenantScope)
		sales.With(read).Get("/", handler.List)
		sales.With(write).Post("/", handler.Create)
		sales.With(read).Get("/{id}", handler.Get)
		sales.With(write).Patch("/{id}", handler.Update)
		sales.With(middleware.RequireRole("admin", "manager")).Delete("/{id}", handler.Delete)
		sales.With(write).Post("/{id}/confirm", handler.Confirm)
		sales.With(middleware.RequireRole("admin", "manager")).Post("/{id}/cancel", handler.Cancel)
	})
}
