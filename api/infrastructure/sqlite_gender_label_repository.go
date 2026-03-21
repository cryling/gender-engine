package infrastructure

import (
	"database/sql"
	"log"

	"github.com/cryling/gender-engine/api/domain"
	_ "github.com/mattn/go-sqlite3"
)

type GenderLabelStorage struct {
	db          *sql.DB
	stmtName    *sql.Stmt
	stmtCountry *sql.Stmt
}

func NewGenderLabelStorage(db *sql.DB) *GenderLabelStorage {
	stmtName, err := db.Prepare("SELECT name, gender FROM gender_labels WHERE name = ? LIMIT 1")
	if err != nil {
		log.Fatalf("Failed to prepare gender_labels statement: %v", err)
	}
	stmtCountry, err := db.Prepare("SELECT name, gender, country, probability FROM gender_country_labels WHERE name = ? AND country = ? ORDER BY probability DESC LIMIT 1")
	if err != nil {
		log.Fatalf("Failed to prepare gender_country_labels statement: %v", err)
	}
	return &GenderLabelStorage{db: db, stmtName: stmtName, stmtCountry: stmtCountry}
}

func (handler *GenderLabelStorage) FindByNameAndCountry(name string, country string) (*domain.GenderCountryLabel, error) {
	row := handler.stmtCountry.QueryRow(name, country)

	label := domain.GenderCountryLabel{}
	err := row.Scan(&label.Name, &label.Gender, &label.Country, &label.Probability)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &domain.NotFoundError{Name: name}
		}
		log.Printf("Unexpected error querying gender_country_labels: %v", err)
		return nil, err
	}

	return &label, nil
}

func (handler *GenderLabelStorage) FindByName(name string) (*domain.GenderLabel, error) {
	row := handler.stmtName.QueryRow(name)

	label := domain.GenderLabel{}
	err := row.Scan(&label.Name, &label.Gender)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &domain.NotFoundError{Name: name}
		}
		log.Printf("Unexpected error querying gender_labels: %v", err)
		return nil, err
	}

	return &label, nil
}
