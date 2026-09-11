package middleware

import (
	"log/slog"
	"net/http"

	"github.com/MaiconGambini/erpGolang/backend/internal/shared/httpx"
)

func Recovery(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("panic recovered", slog.Any("panic", rec), slog.String("path", r.URL.Path))
					httpx.Error(w, "INTERNAL_ERROR", "internal server error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
