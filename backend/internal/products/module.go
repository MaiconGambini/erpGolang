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

	r.Route("/products", func(products chi.Router) {
		products.Use(middleware.AuthJWT(authSvc.JWT()))
		products.Use(middleware.TenantScope)
		products.Get("/low-stock", handler.LowStock)
		products.Get("/", handler.List)
		products.Post("/", handler.Create)
		products.Get("/{id}", handler.Get)
		products.Patch("/{id}", handler.Update)
		products.Delete("/{id}", handler.Delete)
	})
}
