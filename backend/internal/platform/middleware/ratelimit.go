package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/MaiconGambini/erpGolang/backend/internal/shared/httpx"
	goredis "github.com/redis/go-redis/v9"
)

func LoginRateLimit(redis *goredis.Client, maxAttempts int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if redis == nil {
				next.ServeHTTP(w, r)
				return
			}
			key := fmt.Sprintf("ratelimit:login:%s", r.RemoteAddr)
			ctx := r.Context()
			count, err := redis.Incr(ctx, key).Result()
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			if count == 1 {
				_ = redis.Expire(ctx, key, window).Err()
			}
			if count > int64(maxAttempts) {
				httpx.Error(w, "RATE_LIMITED", "too many login attempts", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// ensure context import used in future extensions
var _ = context.Background
