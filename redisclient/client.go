package redisclient

import (
	"fmt"

	"github.com/cryling/gender-engine/config"
	"github.com/redis/go-redis/v9"
)

func RedisClient() *redis.Client {

	appConfig := config.LoadConfig()

	opt, _ := redis.ParseURL(appConfig.REDIS_URL)
	client := redis.NewClient(opt)

	fmt.Println(appConfig.REDIS_URL)
	fmt.Println(client)

	return client
}
