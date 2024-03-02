package infrastructure

import (
	"context"
	"log"

	"github.com/cryling/gender-engine/domain"
	"github.com/redis/go-redis/v9"
)

type RedisHandler struct {
	client  *redis.Client
	context context.Context
}

func NewRedisHandler(client *redis.Client, context context.Context) *RedisHandler {
	return &RedisHandler{client: client, context: context}
}

func (handler *RedisHandler) FindByName(name string) (*domain.GenderLabel, error) {
	result, err := handler.client.Get(handler.context, name).Result()
	if err == redis.Nil {
		log.Printf("Name not found")
		return &domain.GenderLabel{}, err
	}

	label := domain.GenderLabel{
		Name:   name,
		Gender: result,
	}

	return &label, nil
}

func (handler *RedisHandler) InitializeDatabase(data [][]string) error {
	ctx := context.Background()

	pipeline := handler.client.Pipeline()

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

	return nil
}
