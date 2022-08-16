package redis

import (
	"context"
	"github.com/go-redis/redis/v9"
)

var rdb *redis.Client
var ctx = context.Background()

func InitClient() (err error) {
	rdb = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "mymaster",
		SentinelAddrs: []string{"120.48.100.204:27001", "121.196.207.80:27001", "43.142.141.48:27001"},
		DB:            1,
		Password:      "nimasile123",
	})
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		return err
	}
	return nil
}
func HotSearchZadd(field string) error {
	cmd := rdb.ZIncrBy(ctx, "HotSearch", 1, field)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func HotSearchZreduce(key string) error {
	result := rdb.ZIncrBy(ctx, "HotSearch", -1, key)
	if result.Err() != nil {
		return result.Err()
	}
	if result.Val() <= 0 {
		rm := rdb.ZRem(ctx, "HotSearch", key)
		if rm.Err() != nil {
			return rm.Err()
		}
	}
	return nil
}
func Topk(num int64) ([]redis.Z, error) {
	cmd := rdb.ZRevRangeWithScores(ctx, "HotSearch", 0, num-1)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	val := cmd.Val()

	return val, nil
}
