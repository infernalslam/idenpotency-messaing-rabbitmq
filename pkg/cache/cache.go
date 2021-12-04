package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
)

func NewConnection(uri string) (*redis.Client, error) {
	r := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "guest",
		DB:       0,
	})
	_, err := r.Ping(context.Background()).Result()
	return r, err
}

type Cache interface {
	// wip
}

type cache struct {
	c *redis.Client
}

func NewCache(c *redis.Client) Cache {
	return &cache{c}
}
