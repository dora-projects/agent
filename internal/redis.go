package internal

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"sync"
)

var Rdb *redis.Client
var onceRedis sync.Once

func RedisInit() {
	onceRedis.Do(func() {
		host := viper.GetString("REDIS_HOST")
		port := viper.GetString("REDIS_PORT")
		pwd := viper.GetString("REDIS_PASSWORD")

		if host == "" || pwd == "" {
			panic(errors.New("missing redis env host or port"))
		}

		Rdb = redis.NewClient(&redis.Options{
			Addr:     host + ":" + port,
			Password: pwd,
			DB:       0, // use default DB
		})

		ctx := context.Background()

		_, err := Rdb.Ping(ctx).Result()
		if err != nil {
			panic(err)
		}
	})
}
