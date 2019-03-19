package dbclient

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/600lyy/accountservice/model"
	"github.com/go-redis/redis"
)

//IRedisClient exposed
type IRedisClient interface {
	OpenRedisDB()
	SetSession(session *model.Session) error
	GetSession(key string) (session model.Session, err error)
}

//RedisClient uses *redis.client
type RedisClient struct {
	redisClient *redis.Client
}

//OpenRedisDB returns a new redis client connecting to DB 0
func (rc *RedisClient) OpenRedisDB() {
	rc.redisClient = redis.NewClient(&redis.Options{
		Addr:	 	"localhost:6379",
		Password:	"",
		DB:			0,	
	})

	pong, err := rc.redisClient.Ping().Result()
	fmt.Println(pong, err)
}

func (rc *RedisClient) SetSession(session *model.Session) error {
	value, err := json.Marshal(session)
	if err != nil {
		log.Printf("json marshal err,%s", err)
		return err
	}

	err = rc.redisClient.Set(session.SessionID, value, 0).Err()
	if err != nil {
		log.Printf("redis set err,%s", err)
		return err
	}
	return nil
}

func (rc *RedisClient) GetSession(key string) (session model.Session, err error) {
	value, err := rc.redisClient.Get(key).Result()
	if err != nil {
		log.Printf("redis get err,%s", err)
		return
	}
	json.Unmarshal([]byte(value), &session)
	return
}
