package config

import (
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	REDIS_URL string
}

var appConfig AppConfig

func Initialize(env string) {
	err := godotenv.Load("./config/" + env + ".env")
	if err != nil {
		panic(err)
	}

	redisURL, exists := os.LookupEnv("REDIS_URL")
	if !exists {
		panic("REDIS_URL is not set")
	}

	appConfig = AppConfig{
		REDIS_URL: redisURL,
	}
}

func LoadConfig() *AppConfig {
	return &appConfig
}
