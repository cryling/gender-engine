package initializers

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func InitializeRedis(ctx context.Context, client *redis.Client, data map[string]string) {
	pipeline := client.Pipeline()

	i := 0
	for name, gender := range data {
		pipeline.SetNX(ctx, name, gender, 0)

		if i%10000 == 0 {
			pipeline.Exec(ctx)
		}
		i++
	}

	_, err := pipeline.Exec(ctx)
	if err != nil {
		panic(err)
	}

}
