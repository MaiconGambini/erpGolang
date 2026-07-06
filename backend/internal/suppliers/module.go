package suppliers

import (
	"github.com/MaiconGambini/erpGolang/backend/internal/app"
	"github.com/MaiconGambini/erpGolang/backend/internal/auth"
	"github.com/MaiconGambini/erpGolang/backend/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
)

type Module struct{}

func NewModule() Module { return Module{} }

func (Module) Name() string { return "suppliers" }

func (Module) Register(r chi.Router, deps app.Deps) {
	svc := NewService(deps.DB, deps.Audit)
	handler := NewHandler(svc)
	authSvc := auth.NewService(deps.DB, deps.Config)
	read := middleware.RequireRole("admin", "manager", "operator", "viewer")
	write := middleware.RequireRole("admin", "manager", "operator")

	r.Route("/suppliers", func(suppliers chi.Router) {
		suppliers.Use(middleware.AuthJWT(authSvc.JWT()))
		suppliers.Use(middleware.TenantScope)
		suppliers.With(read).Get("/", handler.List)
		suppliers.With(read).Get("/{id}", handler.Get)
		suppliers.With(write).Post("/", handler.Create)
		suppliers.With(write).Patch("/{id}", handler.Update)
		suppliers.With(middleware.RequireRole("admin")).Delete("/{id}", handler.Delete)
	})
}
