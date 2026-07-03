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

	r.Route("/suppliers", func(suppliers chi.Router) {
		suppliers.Use(middleware.AuthJWT(authSvc.JWT()))
		suppliers.Use(middleware.TenantScope)
		suppliers.Get("/", handler.List)
		suppliers.Post("/", handler.Create)
		suppliers.Get("/{id}", handler.Get)
		suppliers.Patch("/{id}", handler.Update)
		suppliers.Delete("/{id}", handler.Delete)
	})
}
