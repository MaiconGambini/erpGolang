package users

import (
	"github.com/MaiconGambini/erpGolang/backend/internal/app"
	"github.com/go-chi/chi/v5"
)

type Module struct{}

func NewModule() Module { return Module{} }

func (Module) Name() string { return "users" }

func (Module) Register(r chi.Router, deps app.Deps) {
	_ = deps
	// MVP: user CRUD deferred; admin routes can be added here
}
