package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/MaiconGambini/erpGolang/backend/internal/platform/jwt"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/authctx"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/httpx"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/tenantctx"
)

func AuthJWT(jwtService jwt.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" || !strings.HasPrefix(header, "Bearer ") {
				httpx.Error(w, "UNAUTHORIZED", "missing or invalid authorization header", http.StatusUnauthorized)
				return
			}
			token := strings.TrimPrefix(header, "Bearer ")
			claims, err := jwtService.Validate(token)
			if err != nil {
				httpx.Error(w, "UNAUTHORIZED", "invalid access token", http.StatusUnauthorized)
				return
			}
			ctx := r.Context()
			ctx = tenantctx.WithTenantID(ctx, claims.TenantID)
			ctx = authctx.WithUser(ctx, authctx.User{ID: claims.UserID, Role: claims.Role})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func TenantScope(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := tenantctx.TenantIDFromContext(r.Context()); err != nil {
			httpx.Error(w, "UNAUTHORIZED", "tenant context required", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func TenantID(ctx context.Context) (string, error) {
	return tenantctx.TenantIDFromContext(ctx)
}
