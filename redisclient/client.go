package redisclient

import (
	"github.com/cryling/gender-engine/config"
	"github.com/redis/go-redis/v9"
)

func CreateClient() *redis.Client {
	appConfig := config.LoadConfig()

	opt, err := redis.ParseURL(appConfig.REDIS_URL)
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opt)

	return client
}
