package main

import (
	"database/sql"
	"log"
	"testing"

	"github.com/cryling/gender-engine/api/infrastructure"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory sqlite database: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS gender_labels (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		gender TEXT NOT NULL
	)`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	return db
}

func seedGenderLabels(db *sql.DB) {
	_, err := db.Exec(`INSERT INTO gender_labels (name, gender) VALUES
		('Sam', 'M'),
		('Jordan', 'F')`)
	if err != nil {
		log.Fatalf("Failed to seed gender labels: %v", err)
	}
}

func TestFindByName(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	storage := infrastructure.NewGenderLabelStorage(db)
	seedGenderLabels(db)

	tests := []struct {
		name           string
		expectedGender string
		expectError    bool
	}{
		{"Sam", "M", false},
		{"Jordan", "F", false},
		{"Unknown", "", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			label, err := storage.FindByName(test.name)
			if test.expectError {
				if err == nil {
					t.Errorf("Expected an error for %v, got nil", test.name)
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error for %v, got: %v", test.name, err)
				}
				if label.Name != test.name {
					t.Errorf("Expected name %v, got %v", test.name, label.Name)
				}
				if label.Gender != test.expectedGender {
					t.Errorf("Expected gender %v, got %v", test.expectedGender, label.Gender)
				}
			}
		})
	}
}
