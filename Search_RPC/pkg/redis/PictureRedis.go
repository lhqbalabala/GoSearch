package redis

import (
	"Search_RPC/pkg/logger"
	"context"
	"github.com/go-redis/redis/v9"
)

var prdb *redis.Client
var cotx = context.Background()

func InitPictureRedisClient() (err error) {
	prdb = redis.NewClient(&redis.Options{
		Addr:     "43.142.141.48:6379",
		Password: "nimasile123", // no password set
		DB:       0,             // use default DB
	})
	_, err = prdb.Ping(cotx).Result()
	if err != nil {
		return err
	}
	return nil
}

func IsUpload(key string) bool {

	_, err := prdb.Get(cotx, key).Result()
	if err == redis.Nil {
		return false
	} else if err != nil {
		logger.Logger.Infoln("redis set err:", err)
		return false
	} else {
		return true
	}
}
func Upload(key string) error {
	var err = prdb.Set(cotx, key, "1", 0).Err()
	if err != nil {
		return err
	}
	return nil
}
