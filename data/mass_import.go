package data

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func MassImport(data [][]string, client *redis.Client) {
	ctx := context.Background()

	pipeline := client.Pipeline()

	for i, line := range data {
		pipeline.SetNX(ctx, line[0], line[1], 0)

		if i%10000 == 0 {
			pipeline.Exec(ctx)
		}
	}

	_, err := pipeline.Exec(ctx)
	if err != nil {
		panic(err)
	}
}
