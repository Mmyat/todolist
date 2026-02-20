package configs

import (
    "github.com/redis/go-redis/v9"
    "sync"
)

var (
    rdb  *redis.Client
    rOnce sync.Once
)

func GetRedis() *redis.Client {
    rOnce.Do(func() {
        rdb = redis.NewClient(&redis.Options{
            Addr: "localhost:6379", // Docker သုံးရင် "redis:6379"
        })
    })
    return rdb
}