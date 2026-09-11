package audit

import (
	"github.com/MaiconGambini/erpGolang/backend/internal/app"
	"github.com/MaiconGambini/erpGolang/backend/internal/auth"
	"github.com/MaiconGambini/erpGolang/backend/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
)

type Module struct{}

func NewModule() Module { return Module{} }

func (Module) Name() string { return "audit" }

func (Module) Register(r chi.Router, deps app.Deps) {
	authSvc := auth.NewService(deps.DB, deps.Config)
	listSvc := NewListService(deps.DB)
	listHandler := NewListHandler(listSvc)

	r.Route("/audit-logs", func(logs chi.Router) {
		logs.Use(middleware.AuthJWT(authSvc.JWT()))
		logs.Use(middleware.TenantScope)
		logs.Use(middleware.RequireRole("admin"))
		logs.Get("/", listHandler.List)
	})
}
