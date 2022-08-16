package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"

	"time"
)

//var rdb *redis.ClusterClient
//var ctx = context.Background()
//
//func InitClient() (err error) {
//	rdb = redis.NewClusterClient(&redis.ClusterOptions{
//		Addrs: []string{
//			"43.142.141.48:7001",
//			"120.48.100.204:7002",
//			"121.196.207.80:7003",
//			"121.196.207.80:7001",
//			"43.142.141.48:7002",
//			"120.48.100.204:7003",
//		},
//		Password: "nimasile123",
//		PoolSize: 200,
//	})
//	_, err = rdb.Ping(ctx).Result()
//	if err != nil {
//		return err
//	}
//	return nil
//}
var rdb *redis.Client
var ctx = context.Background()

func InitClient() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "43.142.141.48:6380",
		Password: "nimasile123", // no password set
		DB:       0,             // use default DB
	})
	_, err = rdb.Ping(cotx).Result()
	if err != nil {
		return err
	}
	return nil
}

func Isfound(key string) bool {

	_, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return false
	} else if err != nil {
		fmt.Println(err)
		return false
	} else {
		return true
	}
}
func Set(key string) error {
	var err = rdb.Set(ctx, key, "1", 1*time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}
