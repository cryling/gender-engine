package initializers

import (
	"context"
	"database/sql"
	"log"
)

func InitializeSqlite(ctx context.Context, db *sql.DB, data map[string]string) {
	createTable(db)
	createIndex(db)

	if alreadySetUp(db) {
		log.Println("Database already set up")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// Prepare a statement for inserting data
	stmt, err := tx.Prepare("INSERT INTO gender_labels(name, gender) VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Iterate over the hashmap and insert each item
	for name, gender := range data {
		_, err = stmt.Exec(name, gender)
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
		"gender" TEXT
	);`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}

func createIndex(db *sql.DB) {
	log.Println("Creating index if not exists")
	createIndexSQL := `CREATE INDEX IF NOT EXISTS idx_gender_labels_name ON gender_labels (name);`
	_, err := db.Exec(createIndexSQL)
	if err != nil {
		log.Fatal(err)
	}
}

func alreadySetUp(db *sql.DB) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM gender_labels").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	return (count == 3491141)
}
