package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/cryling/gender-engine/config"
	"github.com/cryling/gender-engine/domain"
	"github.com/cryling/gender-engine/infrastructure"
	"github.com/cryling/gender-engine/redisclient"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	ginEnv := flag.String("environment", "development", "Specify if running in 'development' or 'production'")
	config.Initialize(*ginEnv)

	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	sqliteRepo := infrastructure.NewSQLiteHandler(db)
	redisRepo := infrastructure.NewRedisHandler(redisclient.CreateClient(), context.Background())

	r := gin.Default()

	r.GET("/api/v1/sqlite_gender", func(c *gin.Context) {
		name := c.Query("name")

		genderFinder := domain.NewGenderFinder(sqliteRepo, name)

		result, err := genderFinder.Find()
		if err == redis.Nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "The specified name is not genderizable BEEP BOP"})
			return
		} else if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%s could be found", name), "gender": &result})
	})

	r.GET("/api/v1/redis_gender", func(c *gin.Context) {
		name := c.Query("name")

		genderFinder := domain.NewGenderFinder(redisRepo, name)

		result, err := genderFinder.Find()
		if err == redis.Nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "The specified name is not genderizable BEEP BOP"})
			return
		} else if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%s could be found", name), "gender": &result})
	})

	r.Run()
}
