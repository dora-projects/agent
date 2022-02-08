package internal

import "github.com/spf13/viper"

type Config struct {
	HttpPort      string
	WebHost       string
	RedisHost     string
	RedisPort     string
	RedisPassword string
}

func GetConfig() Config {
	return Config{
		HttpPort: viper.GetString("port"),

		WebHost: viper.GetString("WEB_HOST"),

		RedisHost:     viper.GetString("REDIS_HOST"),
		RedisPort:     viper.GetString("REDIS_PORT"),
		RedisPassword: viper.GetString("REDIS_PASSWORD"),
	}
}
