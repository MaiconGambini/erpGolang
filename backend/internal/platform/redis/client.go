package redis

import "context"

type Client struct{}

func Open(ctx context.Context, redisURL string) (*Client, error) {
	_ = ctx
	_ = redisURL
	return &Client{}, nil
}
