package cache

import (
	"github.com/dunkbing/tinyimg/tinyimg/config"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func GetRedisClient() *redis.Client {
	if redisClient == nil {
		opt, _ := redis.ParseURL(config.RedisUrl)
		redisClient = redis.NewClient(opt)
	}
	return redisClient
}
