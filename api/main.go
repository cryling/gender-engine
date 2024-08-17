package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cryling/gender-engine/api/domain"
	"github.com/cryling/gender-engine/api/infrastructure"
	"github.com/cryling/gender-engine/api/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	genderLabelRepo := infrastructure.NewGenderLabelStorage(db)

	log.Println("Initialized")

	r := gin.Default()
	r.ForwardedByClientIP = true

	if os.Getenv("RATE_LIMIT_ENABLED") == "true" {
		r.Use(middleware.RateLimitMiddleware())
	}

	r.GET("/api/v1/gender", func(c *gin.Context) {
		name := c.Query("name")
		country := c.Query("country")

		genderFinder := domain.NewGenderFinder(genderLabelRepo, name, country)

		if genderFinder.Name == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
			return
		}

		if country == "" {
			result, err := genderFinder.FindByName()
			if _, ok := err.(*domain.NotFoundError); ok {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
				return
			} else if err != nil {
				panic(err)
			}

			c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%s could be found", name), "result": &result})
			return

		}

		if !domain.ValidCountryCodes()[genderFinder.Country] {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid country code"})
			return
		}

		result, err := genderFinder.FindByNameAndCountry()
		if _, ok := err.(*domain.NotFoundError); ok {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		} else if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%s could be found", name), "result": &result})
	})

	r.Run()
}
