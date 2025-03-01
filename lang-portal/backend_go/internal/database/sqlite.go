package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() (*sql.DB, error) {
	if DB != nil {
		return DB, nil
	}

	// NOTE: Ideally this should be parameterized
	dbPath := filepath.Join("data", "learning.db")
	os.MkdirAll("data", 0755)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	DB = db
	log.Println("Database connection established")
	return DB, nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
