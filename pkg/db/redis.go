package db

import "goIM/pkg/redis"

func NewRedis(addr string) *redis.Conn {
	return redis.New(addr)
}
