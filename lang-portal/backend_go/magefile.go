//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Default target when running mage without arguments
var Default = Build

// Build builds the API binary
func Build() error {
	fmt.Println("Building API...")
	return sh.Run("go", "build", "-o", "bin/api", "cmd/api/main.go")
}

// BuildAll builds the API binary for multiple platforms
func BuildAll() error {
	fmt.Println("Building for multiple platforms...")

	// Create bin directory if it doesn't exist
	if err := os.MkdirAll("bin", 0755); err != nil {
		return err
	}

	// Define target platforms
	targets := []struct {
		os   string
		arch string
	}{
		{"linux", "amd64"},
		{"darwin", "amd64"},
		{"darwin", "arm64"},
		{"windows", "amd64"},
	}

	// Build for each target
	for _, target := range targets {
		fmt.Printf("Building for %s/%s...\n", target.os, target.arch)

		outputFile := filepath.Join("bin", fmt.Sprintf("api-%s-%s", target.os, target.arch))
		if target.os == "windows" {
			outputFile += ".exe"
		}

		env := map[string]string{
			"GOOS":   target.os,
			"GOARCH": target.arch,
		}

		if err := sh.RunWith(env, "go", "build", "-o", outputFile, "cmd/api/main.go"); err != nil {
			return err
		}
	}

	return nil
}

// Run starts the API server
func Run() error {
	fmt.Println("Starting API server...")
	return sh.Run("go", "run", "cmd/api/main.go", "serve")
}

// Test runs all tests
func Test() error {
	fmt.Println("Running tests...")
	return sh.Run("go", "test", "./...")
}

// TestVerbose runs all tests with verbose output
func TestVerbose() error {
	fmt.Println("Running tests with verbose output...")
	return sh.Run("go", "test", "-v", "./...")
}

// TestWords runs only the word endpoint tests
func TestWords() error {
	fmt.Println("Running word endpoint tests...")
	return sh.Run("go", "test", "-v", "./internal/api/handlers/word_test.go")
}

// TestGroups runs only the group endpoint tests
func TestGroups() error {
	fmt.Println("Running group endpoint tests...")
	return sh.Run("go", "test", "-v", "./internal/api/handlers/group_test.go")
}

// TestStudyActivities runs only the study activity endpoint tests
func TestStudyActivities() error {
	fmt.Println("Running study activity endpoint tests...")
	return sh.Run("go", "test", "-v", "./internal/api/handlers/study_activity_test.go")
}

// TestDashboard runs only the dashboard endpoint tests
func TestDashboard() error {
	fmt.Println("Running dashboard endpoint tests...")
	return sh.Run("go", "test", "-v", "./internal/api/handlers/dashboard_test.go")
}

// TestCoverage runs tests with coverage report
func TestCoverage() error {
	fmt.Println("Running tests with coverage report...")

	if err := os.MkdirAll("coverage", 0755); err != nil {
		return err
	}

	if err := sh.Run("go", "test", "-coverprofile=coverage/coverage.out", "./..."); err != nil {
		return err
	}

	return sh.Run("go", "tool", "cover", "-html=coverage/coverage.out", "-o", "coverage/coverage.html")
}

// Clean removes build artifacts
func Clean() error {
	fmt.Println("Cleaning...")
	if err := os.RemoveAll("bin"); err != nil {
		return err
	}
	if err := os.RemoveAll("coverage"); err != nil {
		return err
	}
	return nil
}

// Migrate runs database migrations
func Migrate() error {
	fmt.Println("Running migrations...")
	return sh.Run("go", "run", "cmd/api/main.go", "migrate")
}

// Seed populates the database with seed data
func Seed() error {
	fmt.Println("Seeding database...")
	return sh.Run("go", "run", "cmd/api/main.go", "seed")
}

// Dev runs migrations, seeds the database, and starts the server
func Dev() error {
	mg.SerialDeps(Migrate, Seed)
	return Run()
}

// Lint runs golangci-lint
func Lint() error {
	fmt.Println("Running linter...")

	// Check if golangci-lint is installed
	_, err := exec.LookPath("golangci-lint")
	if err != nil {
		fmt.Println("golangci-lint not found, installing...")
		if runtime.GOOS == "windows" {
			return fmt.Errorf("please install golangci-lint manually on Windows")
		}
		if err := sh.Run("go", "install", "github.com/golangci/golangci-lint/cmd/golangci-lint@latest"); err != nil {
			return err
		}
	}

	return sh.Run("golangci-lint", "run", "./...")
}

// Fmt formats Go code
func Fmt() error {
	fmt.Println("Formatting code...")
	return sh.Run("go", "fmt", "./...")
}

// Tidy runs go mod tidy
func Tidy() error {
	fmt.Println("Tidying dependencies...")
	return sh.Run("go", "mod", "tidy")
}

// Benchmark runs benchmarks
func Benchmark() error {
	fmt.Println("Running benchmarks...")
	return sh.Run("go", "test", "-bench=.", "-benchmem", "./...")
}

// Install installs Mage if it's not already installed
func Install() error {
	fmt.Println("Installing Mage...")

	// Check if mage is already installed
	_, err := exec.LookPath("mage")
	if err == nil {
		fmt.Println("Mage is already installed")
		return nil
	}

	return sh.Run("go", "install", "github.com/magefile/mage@latest")
}

// ResetDB removes the database file, recreates it with migrations, and optionally seeds it
func ResetDB(seed bool) error {
	fmt.Println("Resetting database...")

	// Close any open database connections
	fmt.Println("Closing database connections...")
	CloseDB()

	// Remove the database file
	dbPath := filepath.Join("data", "learning.db")
	fmt.Printf("Removing database file: %s\n", dbPath)
	if err := os.Remove(dbPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove database file: %v", err)
	}

	// Run migrations to recreate the schema
	fmt.Println("Running migrations to recreate schema...")
	if err := Migrate(); err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}

	// Seed the database if requested
	if seed {
		fmt.Println("Seeding database with initial data...")
		if err := Seed(); err != nil {
			return fmt.Errorf("failed to seed database: %v", err)
		}
	}

	fmt.Println("Database reset successfully!")
	return nil
}

// ResetDBClean resets the database without seeding it
func ResetDBClean() error {
	return ResetDB(false)
}

// CloseDB calls the database.CloseDB function
func CloseDB() error {
	fmt.Println("Closing database connections...")
	return sh.Run("go", "run", "cmd/api/main.go", "close-db")
}

// ResetDBWithSeed resets the database and seeds it with initial data
func ResetDBWithSeed() error {
	return ResetDB(true)
}

// Ci runs the CI pipeline (lint, test, build)
func Ci() {
	mg.SerialDeps(Lint, Test, Build)
}
