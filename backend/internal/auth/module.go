package auth

import (
	"github.com/MaiconGambini/erpGolang/backend/internal/app"
	"github.com/MaiconGambini/erpGolang/backend/internal/platform/middleware"
	"github.com/go-chi/chi/v5"
	"time"
)

type Module struct{}

func NewModule() Module { return Module{} }

func (Module) Name() string { return "auth" }

func (Module) Register(r chi.Router, deps app.Deps) {
	svc := NewService(deps.DB, deps.Config)
	handler := NewHandler(svc, deps.Config.BcryptCost)

	r.Route("/auth", func(auth chi.Router) {
		auth.With(middleware.LoginRateLimit(deps.Redis, 5, 15*time.Minute)).Post("/login", handler.Login)
		auth.Post("/refresh", handler.Refresh)
		auth.Post("/logout", handler.Logout)

		auth.Group(func(private chi.Router) {
			private.Use(middleware.AuthJWT(svc.JWT()))
			private.Use(middleware.TenantScope)
			private.Get("/me", handler.Me)
		})
	})
}
