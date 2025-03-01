package database

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// TestDB represents a test database
type TestDB struct {
	DB   *sql.DB
	Path string
}

// NewTestDB creates a new test database
func NewTestDB() (*TestDB, error) {
	// Create a temporary file for the test database
	tempFile, err := ioutil.TempFile("", "test-db-*.sqlite")
	if err != nil {
		return nil, err
	}
	tempFile.Close()

	// Open the database connection
	db, err := sql.Open("sqlite3", tempFile.Name())
	if err != nil {
		os.Remove(tempFile.Name())
		return nil, err
	}

	// Create the test database
	testDB := &TestDB{
		DB:   db,
		Path: tempFile.Name(),
	}

	// Initialize the database schema
	if err := testDB.initSchema(); err != nil {
		testDB.Close()
		return nil, err
	}

	return testDB, nil
}

// initSchema initializes the database schema
func (tdb *TestDB) initSchema() error {
	// Find the project root directory
	rootDir, err := findProjectRoot()
	if err != nil {
		return err
	}

	// Read the schema file from the migrations directory
	schemaPath := filepath.Join(rootDir, "internal", "database", "migrations", "01_initial_schema.sql")
	schema, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return err
	}

	// Execute the schema
	_, err = tdb.DB.Exec(string(schema))
	return err
}

// findProjectRoot attempts to find the project root directory
func findProjectRoot() (string, error) {
	// Start with the current working directory
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Check if we're already at the project root
	if isProjectRoot(dir) {
		return dir, nil
	}

	// Walk up the directory tree until we find the project root
	for {
		parent := filepath.Dir(dir)
		if parent == dir {
			// We've reached the filesystem root without finding the project root
			return "", os.ErrNotExist
		}
		dir = parent
		if isProjectRoot(dir) {
			return dir, nil
		}
	}
}

// isProjectRoot checks if the given directory is the project root
func isProjectRoot(dir string) bool {
	// Check for common indicators of the project root
	indicators := []string{
		"go.mod",
		"internal/database/migrations/01_initial_schema.sql",
	}

	for _, indicator := range indicators {
		if _, err := os.Stat(filepath.Join(dir, indicator)); err == nil {
			return true
		}
	}

	return false
}

// Close closes the database connection and removes the temporary file
func (tdb *TestDB) Close() {
	if tdb.DB != nil {
		tdb.DB.Close()
	}
	if tdb.Path != "" {
		os.Remove(tdb.Path)
	}
}
