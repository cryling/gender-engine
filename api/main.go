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
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db := setupDatabase("./data.db")
	defer db.Close()

	genderLabelRepo := infrastructure.NewGenderLabelStorage(db)

	r := setupRouter(genderLabelRepo)

	if err := r.Run(); err != nil {
		log.Fatalf("Failed to run the server: %v", err)
	}
}

func setupDatabase(dbPath string) *sql.DB {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	log.Println("Database connection established")
	return db
}

func setupRouter(genderLabelRepo domain.GenderLabelRepository) *gin.Engine {
	r := gin.Default()
	r.ForwardedByClientIP = true

	if os.Getenv("RATE_LIMIT_ENABLED") == "true" {
		r.Use(middleware.RateLimitMiddleware())
	}

	r.GET("/api/v1/gender", handleGenderRequest(genderLabelRepo))

	return r
}

func handleGenderRequest(genderLabelRepo domain.GenderLabelRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		country := c.Query("country")

		if name == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
			return
		}

		genderFinder := domain.NewGenderFinder(genderLabelRepo, name, country)

		if country == "" {
			handleFindByName(c, genderFinder)
		} else {
			handleFindByNameAndCountry(c, genderFinder, country)
		}
	}
}

func handleFindByName(c *gin.Context, genderFinder *domain.GenderFinder) {
	result, err := genderFinder.FindByName()
	if err != nil {
		handleError(c, err, genderFinder.Name)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%s could be found", genderFinder.Name),
		"result":  result,
	})
}

func handleFindByNameAndCountry(c *gin.Context, genderFinder *domain.GenderFinder, country string) {
	if !domain.ValidCountryCodes()[genderFinder.Country] {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid country code"})
		return
	}

	result, err := genderFinder.FindByNameAndCountry()
	if err != nil {
		handleError(c, err, genderFinder.Name)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%s could be found in %s", genderFinder.Name, country),
		"result":  result,
	})
}

func handleError(c *gin.Context, err error, name string) {
	if _, ok := err.(*domain.NotFoundError); ok {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
	} else {
		log.Printf("An error occurred while processing the request for %s: %v", name, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}
}
