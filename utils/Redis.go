package utils

import (
	"fmt"
	"gin-spa-app-framework-back/configs"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

func connect() *redis.Client {
	config := configs.SystemConfigs.Redis
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password: config.Password,
		DB:       config.Db,
	})
}

func InitRedis() {
	Rdb = connect()
}
