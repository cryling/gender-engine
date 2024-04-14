package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/cryling/gender-engine/config"
	"github.com/cryling/gender-engine/domain"
	"github.com/cryling/gender-engine/infrastructure"
	"github.com/cryling/gender-engine/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	ginEnv := flag.String("environment", "development", "Specify if running in 'development' or 'production'")
	config.Initialize(*ginEnv)

	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	genderLabelRepo := infrastructure.NewGenderLabelStorage(db)

	log.Println("Initialized")

	r := gin.Default()
	r.ForwardedByClientIP = true
	r.Use(middleware.RateLimitMiddleware())

	r.GET("/api/v1/sqlite_gender", func(c *gin.Context) {
		name := c.Query("name")

		genderFinder := domain.NewGenderFinder(genderLabelRepo, name)

		result, err := genderFinder.Find()
		if _, ok := err.(*domain.NotFoundError); ok {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err})
			return
		} else if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%s could be found", name), "gender": &result})
	})

	r.Run()
}
