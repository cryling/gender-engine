package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/cryling/gender-engine/csv-to-sqlite/initializers"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	ctx := context.Background()

	db, err := sql.Open("sqlite3", "./data/data.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	genderCSV := initializers.InitializeGenderCSV("data/wgnd_2_0_name-gender.csv")
	genderCountryCSV := initializers.InitializeGenderCountryCSV("data/wgnd_2_0_name-gender-code.csv")

	initializers.InitializeSqlite(ctx, db)
	initializers.InitializeGender(ctx, db, *genderCSV)
	initializers.InitializeGenderCountry(ctx, db, *genderCountryCSV)
}
