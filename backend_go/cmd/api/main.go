package main

import (
	"flag"
	"log"
	"os"

	"github.com/andlogreg/free-genai-bootcamp-2025/backend_go/internal/api"
	"github.com/andlogreg/free-genai-bootcamp-2025/backend_go/internal/api/handlers"
	"github.com/andlogreg/free-genai-bootcamp-2025/backend_go/internal/database"
	"github.com/andlogreg/free-genai-bootcamp-2025/backend_go/internal/repository"
	"github.com/andlogreg/free-genai-bootcamp-2025/backend_go/internal/service"
)

func main() {
	// Parse command line arguments
	var command string
	flag.StringVar(&command, "command", "serve", "Command to run (serve, migrate, seed)")
	flag.Parse()

	// If no command provided in args, use the first argument
	if flag.NArg() > 0 {
		command = flag.Arg(0)
	}

	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	// Handle different commands
	switch command {
	case "migrate":
		if err := database.RunMigrations(); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		log.Println("Migrations completed successfully")
		os.Exit(0)

	case "seed":
		if err := database.RunSeed(); err != nil {
			log.Fatalf("Failed to seed database: %v", err)
		}
		log.Println("Database seeded successfully")
		os.Exit(0)

	case "close-db":
		database.CloseDB()
		log.Println("Database connections closed")
		os.Exit(0)

	case "serve":
		// Initialize repositories
		wordRepo := repository.NewWordRepository(db)
		groupRepo := repository.NewGroupRepository(db)
		studyActivityRepo := repository.NewStudyActivityRepository(db)
		studySessionRepo := repository.NewStudySessionRepository(db)

		// Initialize services
		dashboardService := service.NewDashboardService(studySessionRepo, wordRepo, groupRepo)
		studyActivityService := service.NewStudyActivityService(studyActivityRepo, studySessionRepo)
		wordService := service.NewWordService(wordRepo)
		groupService := service.NewGroupService(groupRepo)

		// Initialize handlers
		dashboardHandler := handlers.NewDashboardHandler(dashboardService)
		studyActivityHandler := handlers.NewStudyActivityHandler(studyActivityService)
		wordHandler := handlers.NewWordHandler(wordService)
		groupHandler := handlers.NewGroupHandler(groupService)

		// Setup router
		router := api.SetupRouter(dashboardHandler, studyActivityHandler, wordHandler, groupHandler)

		// Start server
		log.Println("Starting server on :8080")
		if err := router.Run(":8080"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}

	default:
		log.Fatalf("Unknown command: %s", command)
	}
}
