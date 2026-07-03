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

	r.Route("/sales", func(sales chi.Router) {
		sales.Use(middleware.AuthJWT(authSvc.JWT()))
		sales.Use(middleware.TenantScope)
		sales.Get("/", handler.List)
		sales.Post("/", handler.Create)
		sales.Get("/{id}", handler.Get)
		sales.Patch("/{id}", handler.Update)
		sales.Delete("/{id}", handler.Delete)
		sales.Post("/{id}/confirm", handler.Confirm)
		sales.Post("/{id}/cancel", handler.Cancel)
	})
}
