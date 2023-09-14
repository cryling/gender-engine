package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/cryling/gender-engine/config"
	"github.com/cryling/gender-engine/data"
	"github.com/cryling/gender-engine/redisclient"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	ginEnv := flag.String("environment", "development", "Specify if running in 'development' or 'production'")

	r := gin.Default()

	config.Initialize(*ginEnv)

	ctx := context.Background()
	client := redisclient.RedisClient()

	nameData := data.ReadCsvFile("name_gender.csv")

	start := time.Now()
	data.MassImport(nameData, client)
	elapsed := time.Since(start)

	fmt.Println(len(nameData))
	fmt.Println("Data import took ", elapsed)

	r.GET("/api/v1/gender", func(c *gin.Context) {
		name := c.Query("name")

		result, err := client.Get(ctx, name).Result()
		if err == redis.Nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "The specified is not genderizable BEEP BOP"})
			return
		} else if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%s could be found", name), "gender": result})
	})

	r.Run()
}
