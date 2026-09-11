package app

import (
	"log/slog"

	"github.com/MaiconGambini/erpGolang/backend/internal/config"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/audit"
	"github.com/MaiconGambini/erpGolang/backend/internal/platform/validation"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"
)

type Module interface {
	Name() string
	Register(r chi.Router, deps Deps)
}

type Deps struct {
	Config    config.Config
	DB        *pgxpool.Pool
	Redis     *goredis.Client
	Validator validation.Validator
	Audit     audit.Recorder
	Logger    *slog.Logger
}

type Dependencies struct {
	Config    config.Config
	DB        *pgxpool.Pool
	Redis     *goredis.Client
	Validator validation.Validator
	Audit     audit.Recorder
	Logger    *slog.Logger
	Modules   []Module
}

func (d Dependencies) deps() Deps {
	return Deps{
		Config:    d.Config,
		DB:        d.DB,
		Redis:     d.Redis,
		Validator: d.Validator,
		Audit:     d.Audit,
		Logger:    d.Logger,
	}
}
