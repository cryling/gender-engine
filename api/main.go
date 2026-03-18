package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/cryling/gender-engine/api/domain"
	"github.com/cryling/gender-engine/api/infrastructure"
	"github.com/cryling/gender-engine/api/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db := setupDatabase("./data.db")
	defer db.Close()

	genderLabelRepo := infrastructure.NewGenderLabelStorage(db)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/gender", handleGenderRequest(genderLabelRepo))

	handler := middleware.RateLimitMiddleware(mux)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
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

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func handleGenderRequest(genderLabelRepo domain.GenderLabelRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		country := r.URL.Query().Get("country")

		if name == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Name is required"})
			return
		}

		genderFinder := domain.NewGenderFinder(genderLabelRepo, name, country)

		if country == "" {
			handleFindByName(w, genderFinder)
		} else {
			handleFindByNameAndCountry(w, genderFinder, country)
		}
	}
}

func handleFindByName(w http.ResponseWriter, genderFinder *domain.GenderFinder) {
	result, err := genderFinder.FindByName()
	if err != nil {
		handleError(w, err, genderFinder.Name)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"message": fmt.Sprintf("%s could be found", genderFinder.Name),
		"result":  result,
	})
}

func handleFindByNameAndCountry(w http.ResponseWriter, genderFinder *domain.GenderFinder, country string) {
	if !domain.ValidCountryCodes()[genderFinder.Country] {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid country code"})
		return
	}

	result, err := genderFinder.FindByNameAndCountry()
	if err != nil {
		handleError(w, err, genderFinder.Name)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"message": fmt.Sprintf("%s could be found in %s", genderFinder.Name, country),
		"result":  result,
	})
}

func handleError(w http.ResponseWriter, err error, name string) {
	if _, ok := err.(*domain.NotFoundError); ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"message": err.Error()})
	} else {
		log.Printf("An error occurred while processing the request for %s: %v", name, err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
}
