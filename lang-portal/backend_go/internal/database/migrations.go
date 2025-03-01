package database

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"
)

func RunMigrations() error {
	db, err := InitDB()
	if err != nil {
		return fmt.Errorf("failed to initialize database: %v", err)
	}

	// Create migrations table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %v", err)
	}

	// Get list of migration files
	files, err := filepath.Glob("internal/database/migrations/*.sql")
	if err != nil {
		return fmt.Errorf("failed to list migration files: %v", err)
	}

	sort.Strings(files)

	// Get applied migrations
	rows, err := db.Query("SELECT name FROM migrations")
	if err != nil {
		return fmt.Errorf("failed to query migrations: %v", err)
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("failed to scan migration row: %v", err)
		}
		applied[name] = true
	}

	// Apply new migrations
	for _, file := range files {
		name := filepath.Base(file)
		if applied[name] {
			continue
		}

		log.Printf("Applying migration: %s", name)

		content, err := ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %v", name, err)
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %v", err)
		}

		// Split the file into individual statements
		statements := strings.Split(string(content), ";")
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}

			if _, err := tx.Exec(stmt); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to execute migration %s: %v", name, err)
			}
		}

		// Record the migration
		if _, err := tx.Exec("INSERT INTO migrations (name) VALUES (?)", name); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %v", name, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %v", name, err)
		}

		log.Printf("Successfully applied migration: %s", name)
	}

	return nil
}
