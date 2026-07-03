package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/MaiconGambini/erpGolang/backend/internal/shared/authctx"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/tenantctx"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func Logger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(sw, r)
			logger.Info("http request",
				slog.String("request_id", RequestIDFromContext(r.Context())),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", sw.status),
				slog.Int64("duration_ms", time.Since(start).Milliseconds()),
				slog.String("tenant_id", tenantIDFromRequest(r)),
				slog.String("user_id", userIDFromRequest(r)),
			)
		})
	}
}

func tenantIDFromRequest(r *http.Request) string {
	id, err := tenantctx.TenantIDFromContext(r.Context())
	if err != nil {
		return ""
	}
	return id
}

func userIDFromRequest(r *http.Request) string {
	user, err := authctx.UserFromContext(r.Context())
	if err != nil {
		return ""
	}
	return user.ID
}
