package app

import (
	"context"
	"net/http"

	"github.com/MaiconGambini/erpGolang/backend/internal/platform/middleware"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/httpx"
	"github.com/go-chi/chi/v5"
)

func NewRouter(deps Dependencies) http.Handler {
	logger := deps.Logger
	if logger == nil {
		panic("app: logger is required")
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger(logger))
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.CORS(deps.Config.AllowedOrigins))

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		httpx.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if deps.DB != nil {
			if err := deps.DB.Ping(ctx); err != nil {
				httpx.Error(w, "NOT_READY", "database unavailable", http.StatusServiceUnavailable)
				return
			}
		}
		if deps.Redis != nil {
			if err := deps.Redis.Ping(ctx).Err(); err != nil {
				httpx.Error(w, "NOT_READY", "redis unavailable", http.StatusServiceUnavailable)
				return
			}
		}
		httpx.JSON(w, http.StatusOK, map[string]string{"status": "ready"})
	})

	r.Route("/api/v1", func(api chi.Router) {
		api.Get("/", func(w http.ResponseWriter, r *http.Request) {
			httpx.JSON(w, http.StatusOK, map[string]string{
				"name":    "goERP API",
				"version": "v1",
			})
		})

		moduleDeps := deps.deps()
		for _, mod := range deps.Modules {
			mod.Register(api, moduleDeps)
		}
	})

	return r
}

func PingDB(ctx context.Context, deps Dependencies) error {
	if deps.DB == nil {
		return nil
	}
	return deps.DB.Ping(ctx)
}
