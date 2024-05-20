package initializers

import (
	"context"
	"database/sql"
	"log"
)

func InitializeSqlite(ctx context.Context, db *sql.DB, data []GenderData) {
	createTable(db)
	createIndex(db)

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// Prepare a statement for inserting data
	stmt, err := tx.Prepare("INSERT INTO gender_labels(name, gender, country, probability) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Iterate over the list and insert each item
	for _, element := range data {
		_, err = stmt.Exec(element.Name, element.Gender, element.Code, element.Probability)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

}

func createTable(db *sql.DB) {
	log.Println("Creating table if not exists")
	createTableSQL := `CREATE TABLE IF NOT EXISTS gender_labels (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,     
		"name" TEXT,
		"country" TEXT,
		"gender" TEXT,
		"probability" REAL
	);`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}

func createIndex(db *sql.DB) {
	log.Println("Creating index if not exists")
	createIndexSQL := `CREATE INDEX IF NOT EXISTS idx_gender_labels_name ON gender_labels (name, country, probability);`
	_, err := db.Exec(createIndexSQL)
	if err != nil {
		log.Fatal(err)
	}
}
