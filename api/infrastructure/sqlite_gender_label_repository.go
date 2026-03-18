package infrastructure

import (
	"database/sql"
	"log"

	"github.com/cryling/gender-engine/api/domain"
	_ "github.com/mattn/go-sqlite3"
)

type GenderLabelStorage struct {
	DB *sql.DB
}

func NewGenderLabelStorage(db *sql.DB) *GenderLabelStorage {
	return &GenderLabelStorage{DB: db}
}

func (handler *GenderLabelStorage) FindByNameAndCountry(name string, country string) (*domain.GenderCountryLabel, error) {
	row := handler.DB.QueryRow(
		"SELECT name, gender, country, probability FROM gender_country_labels WHERE name = ? AND country = ? ORDER BY probability DESC LIMIT 1",
		name,
		country,
	)

	label := domain.GenderCountryLabel{}
	err := row.Scan(&label.Name, &label.Gender, &label.Country, &label.Probability)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &domain.NotFoundError{Name: name}
		}
		log.Printf("Unexpected error querying gender_country_labels for %s: %v", name, err)
		return nil, err
	}

	return &label, nil
}

func (handler *GenderLabelStorage) FindByName(name string) (*domain.GenderLabel, error) {
	row := handler.DB.QueryRow(
		"SELECT name, gender FROM gender_labels WHERE name = ?",
		name,
	)

	label := domain.GenderLabel{}
	err := row.Scan(&label.Name, &label.Gender)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &domain.NotFoundError{Name: name}
		}
		log.Printf("Unexpected error querying gender_labels for %s: %v", name, err)
		return nil, err
	}

	return &label, nil
}
