package dbclient

import (
	"fmt"

	"github.com/go-redis/redis"
)

//IRedisClient exposed
type IRedisClient interface {
	OpenRedisDB()
	Set(key, value string) bool
	Get(key string) string
}

//RedisClient uses *redis.client
type RedisClient struct {
	redisClient *redis.Client
}

func (rc *RedisClient) OpenRedisDB() {
	rc.redisClient = redis.NewClient(&redis.Options{
		Addr:	 	"localhost:6379",
		Password:	"",
		DB:			0,	
	})

	pong, err := rc.redisClient.Ping().Result()
	fmt.Println(pong, err)
}
