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

	r.Route("/customers", func(customers chi.Router) {
		customers.Use(middleware.AuthJWT(authSvc.JWT()))
		customers.Use(middleware.TenantScope)
		customers.Get("/", handler.List)
		customers.Post("/", handler.Create)
		customers.Get("/{id}", handler.Get)
		customers.Patch("/{id}", handler.Update)
		customers.Delete("/{id}", handler.Delete)
	})
}
