package audit

import (
	"github.com/MaiconGambini/erpGolang/backend/internal/app"
	"github.com/go-chi/chi/v5"
)

type Module struct{}

func NewModule() Module { return Module{} }

func (Module) Name() string { return "audit" }

func (Module) Register(r chi.Router, deps app.Deps) {
	_ = deps
	// Audit is write-side only via audit.Recorder; no public routes in MVP
}
