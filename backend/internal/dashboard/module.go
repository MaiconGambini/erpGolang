package dashboard

import (
	"github.com/MaiconGambini/erpGolang/backend/internal/app"
	"github.com/MaiconGambini/erpGolang/backend/internal/auth"
	"github.com/MaiconGambini/erpGolang/backend/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
)

type Module struct{}

func NewModule() Module { return Module{} }

func (Module) Name() string { return "dashboard" }

func (Module) Register(r chi.Router, deps app.Deps) {
	svc := NewService(deps.DB)
	handler := NewHandler(svc)
	authSvc := auth.NewService(deps.DB, deps.Config)

	r.Route("/dashboard", func(dashboard chi.Router) {
		dashboard.Use(middleware.AuthJWT(authSvc.JWT()))
		dashboard.Use(middleware.TenantScope)
		dashboard.Get("/summary", handler.Summary)
	})
}
