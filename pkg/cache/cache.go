package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisConf struct {
	Addr     string
	Password string
	DB       int
}

func NewConnection(rConf RedisConf) (*redis.Client, error) {
	r := redis.NewClient(&redis.Options{
		Addr:     rConf.Addr,
		Password: rConf.Password,
		DB:       rConf.DB,
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
