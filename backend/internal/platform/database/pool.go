package database

import "context"

type Pool struct{}

func Open(ctx context.Context, databaseURL string) (*Pool, error) {
	_ = ctx
	_ = databaseURL
	return &Pool{}, nil
}
