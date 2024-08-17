package initializers

import (
	"context"
	"database/sql"
	"log"
)

func InitializeSqlite(ctx context.Context, db *sql.DB) {
	createTablesAndIndexes(db)
}

func InitializeGenderCountry(ctx context.Context, db *sql.DB, genderCountryData []GenderCountryData) {
	executeTransaction(ctx, db, func(tx *sql.Tx) error {
		return insertGenderCountryData(tx, genderCountryData)
	})
}

func InitializeGender(ctx context.Context, db *sql.DB, genderData []GenderData) {
	executeTransaction(ctx, db, func(tx *sql.Tx) error {
		return insertGenderData(tx, genderData)
	})
}

func createTablesAndIndexes(db *sql.DB) {
	execStatements(db, []string{
		`CREATE TABLE IF NOT EXISTS gender_labels (
			"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,     
			"name" TEXT,
			"gender" TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS gender_country_labels (
			"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,     
			"name" TEXT,
			"country" TEXT,
			"gender" TEXT,
			"probability" REAL
		);`,
		`CREATE INDEX IF NOT EXISTS idx_gender_labels_name ON gender_labels (name);`,
		`CREATE INDEX IF NOT EXISTS idx_gender_country_labels_name ON gender_country_labels (name, country, probability);`,
	})
}

func execStatements(db *sql.DB, statements []string) {
	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			log.Fatalf("Failed to execute statement: %v", err)
		}
	}
}

func executeTransaction(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}

	if err := fn(tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Fatalf("Failed to rollback transaction: %v", rollbackErr)
		}
		log.Fatalf("Failed to execute transaction: %v", err)
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}
}

func insertGenderData(tx *sql.Tx, data []GenderData) error {
	stmt, err := tx.Prepare("INSERT INTO gender_labels(name, gender) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, element := range data {
		if _, err = stmt.Exec(element.Name, element.Gender); err != nil {
			return err
		}
	}
	return nil
}

func insertGenderCountryData(tx *sql.Tx, data []GenderCountryData) error {
	stmt, err := tx.Prepare("INSERT INTO gender_country_labels(name, gender, country, probability) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, element := range data {
		if _, err = stmt.Exec(element.Name, element.Gender, element.Code, element.Probability); err != nil {
			return err
		}
	}
	return nil
}
