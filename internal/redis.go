package internal

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"sync"
)

var Rdb *redis.Client
var runOnce sync.Once

func RedisInit() {
	runOnce.Do(func() {
		conf := GetConfig()

		host := conf.RedisHost
		port := conf.RedisPort
		pwd := conf.RedisPassword

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

type SiteConfig struct {
	Release  string            `json:"release"`
	Filepath string            `json:"filepath"`
	Index    string            `json:"index"`
	Proxy    map[string]string `json:"proxy"`
}

func ParseSiteConfig(data []byte) (*SiteConfig, error) {
	var conf SiteConfig
	err := json.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, err
}
