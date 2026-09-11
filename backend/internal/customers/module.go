package customers

import (
	"github.com/MaiconGambini/erpGolang/backend/internal/app"
	"github.com/MaiconGambini/erpGolang/backend/internal/auth"
	"github.com/MaiconGambini/erpGolang/backend/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
)

type Module struct{}

func NewModule() Module { return Module{} }

func (Module) Name() string { return "customers" }

func (Module) Register(r chi.Router, deps app.Deps) {
	svc := NewService(deps.DB, deps.Audit)
	handler := NewHandler(svc)
	authSvc := auth.NewService(deps.DB, deps.Config)
	read := middleware.RequireRole("admin", "manager", "operator", "viewer")
	write := middleware.RequireRole("admin", "manager", "operator")

	r.Route("/customers", func(customers chi.Router) {
		customers.Use(middleware.AuthJWT(authSvc.JWT()))
		customers.Use(middleware.TenantScope)
		customers.With(read).Get("/", handler.List)
		customers.With(read).Get("/{id}", handler.Get)
		customers.With(write).Post("/", handler.Create)
		customers.With(write).Patch("/{id}", handler.Update)
		customers.With(middleware.RequireRole("admin")).Delete("/{id}", handler.Delete)
	})
}
