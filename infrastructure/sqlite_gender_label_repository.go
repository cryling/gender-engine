package infrastructure

import (
	"database/sql"
	"log"

	"github.com/cryling/gender-engine/domain"
	_ "github.com/mattn/go-sqlite3"
)

type GenderLabelStorage struct {
	DB *sql.DB
}

func NewGenderLabelStorage(db *sql.DB) *GenderLabelStorage {
	return &GenderLabelStorage{DB: db}
}

func (handler *GenderLabelStorage) FindByName(name string) (*domain.GenderLabel, error) {
	row := handler.DB.QueryRow("SELECT * FROM gender_labels WHERE name = ? LIMIT 1", name)

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
