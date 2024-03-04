package infrastructure

import (
	"database/sql"
	"log"

	"github.com/cryling/gender-engine/domain"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteHandler struct {
	DB *sql.DB
}

func NewSQLiteHandler(db *sql.DB) *SQLiteHandler {
	return &SQLiteHandler{DB: db}
}

func (handler *SQLiteHandler) FindByName(name string) (*domain.GenderLabel, error) {
	row := handler.DB.QueryRow("SELECT * FROM data WHERE name = ? LIMIT 1", name)

	label := domain.GenderLabel{}

	var err error

	if err = row.Scan(&label.Name, &label.Gender); err == sql.ErrNoRows {
		log.Printf("Name not found")
		return &domain.GenderLabel{}, err
	}

	return &label, nil
}
