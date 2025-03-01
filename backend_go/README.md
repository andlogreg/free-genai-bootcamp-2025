# Language Learning Portal Backend

This is the backend server for the Language Learning Portal, built with Go and Gin framework.

## Project Structure

```
backend_go/
├── cmd/                      # Application entry points
│   └── api/                 # Main API server
│       └── main.go         # Server initialization and configuration
├── internal/                # Private application code
│   ├── api/                # API layer
│   │   ├── handlers/      # HTTP request handlers
│   │   ├── middleware/    # HTTP middleware
│   │   └── router.go      # Route definitions
│   ├── models/            # Database models
│   ├── repository/        # Database operations
│   ├── service/          # Business logic
│   └── database/         # Database configuration and migrations
│       ├── migrations/   # SQL migration files
│       └── sqlite.go     # SQLite connection and configuration
├── pkg/                   # Public library code
│   └── utils/            # Shared utilities
├── seeds/                 # Seed data JSON files
├── data/                  # Database files
├── bin/                   # Compiled binaries
├── magefile.go            # Mage build tasks
├── mage.sh                # Mage wrapper script
├── go.mod                 # Go module definition
├── go.sum                 # Go module checksums
└── README.md              # This file
```

## Requirements

- Go 1.21 or higher
- SQLite3
- Mage (optional, for development tasks)

## Setup

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Run migrations:
   ```bash
   go run cmd/api/main.go migrate
   ```
   This will create all necessary database tables and indexes.

3. Import seed data:
   ```bash
   go run cmd/api/main.go seed
   ```
   This will populate the database with:
   - Sample study activities (Flashcards, Word Matching, Writing Practice)
   - Word groups (Basic Greetings, Numbers 1-10, Common Colors)
   - Portuguese-English word pairs for each group
   - A sample study session with word reviews

4. Start the server:
   ```bash
   go run cmd/api/main.go serve
   ```
   The server will start on port 8080 (http://localhost:8080).

## Development

The project uses Mage for common development tasks. You can run Mage targets using:

```bash
# Using the wrapper script
./mage.sh <target>

# Or directly with go run
go run -tags=mage magefile.go <target>
```

### Basic Targets

- `build` - Build the API binary (default target)
- `run` - Start the API server
- `test` - Run all tests
- `clean` - Remove build artifacts
- `migrate` - Run database migrations
- `seed` - Seed the database with sample data
- `dev` - Run migrations, seed the database, and start the server
- `resetdbclean` - Reset the database by removing the database file and recreating the schema (without seeding)
- `resetdbwithseed` - Reset the database and seed it with initial data

### Testing Targets

- `test` - Run all tests
- `testWords` - Run only word endpoint tests
- `testGroups` - Run only group endpoint tests
- `testStudyActivities` - Run only study activity endpoint tests
- `testDashboard` - Run only dashboard endpoint tests
- `testIntegration` - Run integration tests
- `testUnit` - Run unit tests (skips integration tests)
- `testVerbose` - Run tests with verbose output
- `testCoverage` - Run tests with coverage report (generates HTML report in coverage directory)

### Advanced Targets

- `buildall` - Build the API binary for multiple platforms (Linux, macOS, Windows)
- `lint` - Run golangci-lint (installs it if not found)
- `fmt` - Format Go code
- `tidy` - Run go mod tidy
- `benchmark` - Run benchmarks
- `ci` - Run the CI pipeline (lint, test, build)

To see all available targets:
```bash
./mage.sh -l
```

## API Documentation

The API provides the following endpoints:

### Dashboard
- `GET /api/dashboard/last_study_session` - Get the most recent study session
- `GET /api/dashboard/study_progress` - Get study progress statistics
- `GET /api/dashboard/quick-stats` - Get quick statistics about words, groups, and study sessions

### Study Activities
- `GET /api/study_activities` - List all study activities
- `GET /api/study_activities/:id` - Get a specific study activity
- `GET /api/study_activities/:id/study_sessions` - Get study sessions for a specific activity
- `POST /api/study_activities` - Create a new study activity

### Words
- `GET /api/words` - List all words
- `GET /api/words/:id` - Get a specific word

### Groups
- `GET /api/groups` - List all groups
- `GET /api/groups/:id` - Get a specific group
- `GET /api/groups/:id/words` - Get words in a specific group
- `GET /api/groups/:id/study_sessions` - Get study sessions for a specific group
- `POST /api/groups/:id/words` - Add words to a group (expects an array of word IDs)
- `DELETE /api/groups/:id/words/:word_id` - Remove a word from a group

### Study Sessions
- `GET /api/study_sessions` - List all study sessions
- `GET /api/study_sessions/:id` - Get a specific study session
- `GET /api/study_sessions/:id/words` - Get words reviewed in a specific study session
- `POST /api/study_sessions` - Create a new study session

## Pagination

All list endpoints support pagination with the following query parameters:

- `page` - The page number (default: 1)
- `page_size` - The number of items per page (default: 10, max: 100)

Example request:
```
GET /api/words?page=2&page_size=20
```

Example response:
```json
{
  "items": [
    {
      "id": 21,
      "portuguese": "casa",
      "english": "house",
      "created_at": "2023-01-01T12:00:00Z",
      "correct_count": 5,
      "wrong_count": 2
    },
    // ... more items
  ],
  "pagination": {
    "current_page": 2,
    "total_pages": 5,
    "total_items": 100,
    "items_per_page": 20
  }
}
```

## Project Structure Explanation

### Core Components

- `internal/`: Contains private code that's specific to your application and not meant to be imported by other projects.
   - `api/`: Contains all API-related code
      - `handlers/`: HTTP request handlers that process incoming requests and return responses
      - `middleware/`: Functions that run before or after request handlers (for authentication, logging, etc.)
      - `router.go`: Defines your API routes and connects them to handlers
   - `models/`: Data structures that represent your database tables or API resources
   - `repository/`: Code that handles database operations (creating, reading, updating, deleting data)
   - `service/`: Contains your business logic, sitting between handlers and repositories
   - `database/`: Database configuration and setup
      - `migrations/`: SQL files that define database schema changes
      - `sqlite.go`: Code to connect to and configure your SQLite database
- `pkg/`: Contains public code that could potentially be used by other projects
   - `utils/`: Shared utility functions that might be used across your application
- `seeds/`: Contains JSON files with initial data to populate your database for testing or initial setup

### Architecture Explained

The application follows a layered architecture:

1. **API Layer** (`handlers/`): Handles HTTP requests and responses
2. **Service Layer** (`service/`): Contains business logic
3. **Repository Layer** (`repository/`): Manages data access
4. **Model Layer** (`models/`): Defines data structures

This separation of concerns makes the code more maintainable and testable.

#### Going deeper: Further explanation with analogies

- `api/` - **The Front Desk**
   - Think of the api/ directory as the front desk of a hotel. It's where all external interactions happen.

- `handlers/` - **The Receptionists**
   - **What they do**: Process incoming HTTP requests and return appropriate responses
   - **Example**: A handler for /users/123 might fetch user data and return it as JSON
   - **Analogy**: Like hotel receptionists who take your requests ("I need a room" or "Where's the gym?") and either handle them directly or direct them to the right department

- `middleware/` - **The Security Guards & Assistants**
   - **What they do**: Intercept requests before they reach handlers or after handlers process them
   - **Example**: Authentication middleware checks if a request has a valid token before allowing access
   - **Analogy**: Like security guards who check your key card before letting you into certain areas, or assistants who prepare rooms before guests arrive

- `router.go` - The Directory
   - **What it does**: Maps URLs to the appropriate handlers
   - **Example**: Defines that `GET /users/{id}` should be handled by the `GetUserHandler` function
   - **Analogy**: Like the hotel directory that tells you which floor and room number to go to for different services

- `models/` - **The Blueprints** - Think of models/ as the blueprints or templates that define what your data looks like.

   - **What they do**: Define data structures that represent your database tables or API resources
   - **Example**: A `User` struct with fields like `ID`, `Name`, `Email`, etc.
   - **Analogy**: Like the standard room layouts in a hotel - they define what a "Deluxe Suite" or "Standard Room" contains

- `service/` - **The Management** - The `service/` directory contains your business logic - the rules and processes of your application.

   - **What it does**: Implements the core functionality of your application
   - **Example**: A UserService might contain functions for creating users, validating emails, or handling password resets
   - **Analogy**: Like hotel management that makes decisions based on business rules - "If a guest stays 7+ nights, offer a discount" or "If the hotel is at 90% capacity, raise rates"
- `repository/` - **The Record Keepers** - The `repository/` handles all database interactions, keeping your data access code separate from business logic.

   - **What it does**: Provides methods to create, read, update, and delete data from your database
   - **Example**: A UserRepository might have methods like GetUserByID(), CreateUser(), or UpdateUserEmail()
   - **Analogy**: Like the hotel's record-keeping department that manages guest records, room availability, and booking history

**Real-world Flow Example**

For a hypothetical "create user" feature:

1. **Router**: Routes `POST /users` request to the CreateUserHandler
2. **Middleware**: Validates the API key and logs the request
3. **Handler**: Parses the JSON request body into a user model
4. **Service**: Validates the user data, hashes the password, and calls the repository
5. **Repository**: Inserts the user data into the database
6. **Model**: Defines what user data looks like throughout this process