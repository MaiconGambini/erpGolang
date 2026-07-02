package authctx

import (
	"context"
	"errors"
)

type User struct {
	ID   string
	Role string
}

type contextKey struct{}

func WithUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, contextKey{}, user)
}

func UserFromContext(ctx context.Context) (User, error) {
	value, ok := ctx.Value(contextKey{}).(User)
	if !ok || value.ID == "" {
		return User{}, errors.New("user missing from context")
	}
	return value, nil
}
