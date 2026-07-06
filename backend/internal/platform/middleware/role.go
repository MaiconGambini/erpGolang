package middleware

import (
	"net/http"
	"slices"

	"github.com/MaiconGambini/erpGolang/backend/internal/shared/authctx"
	"github.com/MaiconGambini/erpGolang/backend/internal/shared/httpx"
)

func RequireRole(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := authctx.UserFromContext(r.Context())
			if err != nil {
				httpx.Error(w, "UNAUTHORIZED", "authentication required", http.StatusUnauthorized)
				return
			}
			if _, ok := allowed[user.Role]; !ok {
				httpx.Error(w, "FORBIDDEN", "insufficient permissions", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func HasRole(role string, roles ...string) bool {
	return slices.Contains(roles, role)
}
