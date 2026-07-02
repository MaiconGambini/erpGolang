package tenantctx

import (
	"context"
	"errors"
)

type contextKey struct{}

func WithTenantID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, contextKey{}, tenantID)
}

func TenantIDFromContext(ctx context.Context) (string, error) {
	value, ok := ctx.Value(contextKey{}).(string)
	if !ok || value == "" {
		return "", errors.New("tenant id missing from context")
	}
	return value, nil
}
