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

	appConfig = AppConfig{
		REDIS_URL: os.Getenv("REDIS_URL"),
	}
}

func LoadConfig() *AppConfig {
	return &appConfig
}
