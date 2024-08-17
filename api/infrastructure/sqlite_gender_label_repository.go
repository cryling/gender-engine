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
		"SELECT * FROM gender_country_labels WHERE name = ? AND country = ? ORDER BY probability DESC LIMIT 1",
		name,
		country,
	)

	label := domain.GenderCountryLabel{}

	var id int
	err := row.Scan(&id, &label.Name, &label.Gender, &label.Country, &label.Probability)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			log.Printf("Name not found")
		default:
			log.Println("Unexpected error: ", err)
		}
		return &domain.GenderCountryLabel{}, &domain.NotFoundError{Name: name}
	}

	return &label, nil
}

func (handler *GenderLabelStorage) FindByName(name string) (*domain.GenderLabel, error) {
	row := handler.DB.QueryRow(
		"SELECT * FROM gender_labels WHERE name = ?",
		name,
	)

	label := domain.GenderLabel{}

	var id int
	err := row.Scan(&id, &label.Name, &label.Gender)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			log.Printf("Name not found")
		default:
			log.Println("Unexpected error: ", err)
		}
		return &domain.GenderLabel{}, &domain.NotFoundError{Name: name}
	}

	return &label, nil
}
