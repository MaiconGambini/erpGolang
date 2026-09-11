package products

import (
	"github.com/MaiconGambini/erpGolang/backend/internal/app"
	"github.com/MaiconGambini/erpGolang/backend/internal/auth"
	"github.com/MaiconGambini/erpGolang/backend/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
)

type Module struct{}

func NewModule() Module { return Module{} }

func (Module) Name() string { return "products" }

func (Module) Register(r chi.Router, deps app.Deps) {
	svc := NewService(deps.DB, deps.Audit)
	handler := NewHandler(svc)
	authSvc := auth.NewService(deps.DB, deps.Config)
	read := middleware.RequireRole("admin", "manager", "operator", "viewer")
	write := middleware.RequireRole("admin", "manager", "operator")

	r.Route("/products", func(products chi.Router) {
		products.Use(middleware.AuthJWT(authSvc.JWT()))
		products.Use(middleware.TenantScope)
		products.With(read).Get("/low-stock", handler.LowStock)
		products.With(read).Get("/", handler.List)
		products.With(read).Get("/{id}", handler.Get)
		products.With(write).Post("/", handler.Create)
		products.With(write).Patch("/{id}", handler.Update)
		products.With(middleware.RequireRole("admin")).Delete("/{id}", handler.Delete)
	})
}
