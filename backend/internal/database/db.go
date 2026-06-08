package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(dbUrl string) *sql.DB {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Failed to open database connection: ", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database: ", err)
	}

	log.Println("Successfully connected to the database")
	return db
}
