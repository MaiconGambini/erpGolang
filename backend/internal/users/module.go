package users

import (
	"github.com/MaiconGambini/erpGolang/backend/internal/app"
	"github.com/MaiconGambini/erpGolang/backend/internal/auth"
	"github.com/MaiconGambini/erpGolang/backend/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
)

type Module struct{}

func NewModule() Module { return Module{} }

func (Module) Name() string { return "users" }

func (Module) Register(r chi.Router, deps app.Deps) {
	svc := NewService(deps.DB)
	handler := NewHandler(svc)
	authSvc := auth.NewService(deps.DB, deps.Config)

	r.Route("/users", func(users chi.Router) {
		users.Use(middleware.AuthJWT(authSvc.JWT()))
		users.Use(middleware.TenantScope)
		users.Use(middleware.RequireRole("admin"))
		users.Get("/", handler.List)
		users.Get("/{id}", handler.Get)
		users.Patch("/{id}", handler.Update)
	})
}
