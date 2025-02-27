package handlers_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/andlogreg/free-genai-bootcamp-2025/backend_go/internal/api"
	"github.com/andlogreg/free-genai-bootcamp-2025/backend_go/internal/api/handlers"
	"github.com/andlogreg/free-genai-bootcamp-2025/backend_go/internal/database"
	"github.com/andlogreg/free-genai-bootcamp-2025/backend_go/internal/models"
	"github.com/andlogreg/free-genai-bootcamp-2025/backend_go/internal/repository"
	"github.com/andlogreg/free-genai-bootcamp-2025/backend_go/internal/service"
	"github.com/andlogreg/free-genai-bootcamp-2025/backend_go/internal/testutil"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// DashboardHandlerTestSuite is a test suite for the dashboard handlers
type DashboardHandlerTestSuite struct {
	suite.Suite
	router              *gin.Engine
	db                  *database.TestDB
	dashboardHandler    *handlers.DashboardHandler
	testWords           []*models.Word
	testGroups          []*models.Group
	testStudyActivities []*models.StudyActivity
	testStudySessions   []*models.StudySession
	testWordReviewItems []*models.WordReviewItem
}

// SetupSuite sets up the test suite
func (suite *DashboardHandlerTestSuite) SetupSuite() {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a temporary test database
	var err error
	suite.db, err = database.NewTestDB()
	if err != nil {
		suite.T().Fatalf("Failed to create test database: %v", err)
	}

	// Initialize repositories
	wordRepo := repository.NewWordRepository(suite.db.DB)
	groupRepo := repository.NewGroupRepository(suite.db.DB)
	studySessionRepo := repository.NewStudySessionRepository(suite.db.DB)
	studyActivityRepo := repository.NewStudyActivityRepository(suite.db.DB)

	// Initialize services
	wordService := service.NewWordService(wordRepo)
	groupService := service.NewGroupService(groupRepo)
	studyActivityService := service.NewStudyActivityService(studyActivityRepo, studySessionRepo)
	dashboardService := service.NewDashboardService(studySessionRepo, wordRepo, groupRepo)

	// Initialize handlers
	wordHandler := handlers.NewWordHandler(wordService)
	groupHandler := handlers.NewGroupHandler(groupService)
	studyActivityHandler := handlers.NewStudyActivityHandler(studyActivityService)
	suite.dashboardHandler = handlers.NewDashboardHandler(dashboardService)

	// Setup router
	suite.router = api.SetupRouter(
		suite.dashboardHandler,
		studyActivityHandler,
		wordHandler,
		groupHandler,
	)
}

// TearDownSuite tears down the test suite
func (suite *DashboardHandlerTestSuite) TearDownSuite() {
	// Close and remove the test database
	if suite.db != nil {
		suite.db.Close()
	}
}

// SetupTest sets up each test
func (suite *DashboardHandlerTestSuite) SetupTest() {
	// Clear any existing test data
	suite.clearTestData()

	// Seed test data for this test
	suite.seedTestData()
}

// TearDownTest cleans up after each test
func (suite *DashboardHandlerTestSuite) TearDownTest() {
	// Clear test data
	suite.clearTestData()
}

// clearTestData removes all test data from the database
func (suite *DashboardHandlerTestSuite) clearTestData() {
	// Clear the test slices
	suite.testWords = nil
	suite.testGroups = nil
	suite.testStudyActivities = nil
	suite.testStudySessions = nil
	suite.testWordReviewItems = nil

	// Use a transaction to handle errors gracefully
	tx, err := suite.db.DB.Begin()
	if err != nil {
		suite.T().Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Check if the word_review_items table exists
	var tableExists int
	err = suite.db.DB.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='word_review_items'").Scan(&tableExists)
	if err != nil {
		suite.T().Fatalf("Failed to check if word_review_items table exists: %v", err)
	}

	if tableExists > 0 {
		// Table exists, so delete the data
		_, err = tx.Exec("DELETE FROM word_review_items")
		if err != nil {
			suite.T().Fatalf("Failed to clear word_review_items data: %v", err)
		}
	}

	// Check if the study_sessions table exists
	err = suite.db.DB.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='study_sessions'").Scan(&tableExists)
	if err != nil {
		suite.T().Fatalf("Failed to check if study_sessions table exists: %v", err)
	}

	if tableExists > 0 {
		// Table exists, so delete the data
		_, err = tx.Exec("DELETE FROM study_sessions")
		if err != nil {
			suite.T().Fatalf("Failed to clear study_sessions data: %v", err)
		}
	}

	// Check if the study_activities table exists
	err = suite.db.DB.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='study_activities'").Scan(&tableExists)
	if err != nil {
		suite.T().Fatalf("Failed to check if study_activities table exists: %v", err)
	}

	if tableExists > 0 {
		// Table exists, so delete the data
		_, err = tx.Exec("DELETE FROM study_activities")
		if err != nil {
			suite.T().Fatalf("Failed to clear study_activities data: %v", err)
		}
	}

	// Check if the words_groups table exists
	err = suite.db.DB.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='words_groups'").Scan(&tableExists)
	if err != nil {
		suite.T().Fatalf("Failed to check if words_groups table exists: %v", err)
	}

	if tableExists > 0 {
		// Table exists, so delete the data
		_, err = tx.Exec("DELETE FROM words_groups")
		if err != nil {
			suite.T().Fatalf("Failed to clear words_groups data: %v", err)
		}
	}

	// Check if the groups table exists
	err = suite.db.DB.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='groups'").Scan(&tableExists)
	if err != nil {
		suite.T().Fatalf("Failed to check if groups table exists: %v", err)
	}

	if tableExists > 0 {
		// Table exists, so delete the data
		_, err = tx.Exec("DELETE FROM groups")
		if err != nil {
			suite.T().Fatalf("Failed to clear groups data: %v", err)
		}
	}

	// Check if the words table exists
	err = suite.db.DB.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='words'").Scan(&tableExists)
	if err != nil {
		suite.T().Fatalf("Failed to check if words table exists: %v", err)
	}

	if tableExists > 0 {
		// Table exists, so delete the data
		_, err = tx.Exec("DELETE FROM words")
		if err != nil {
			suite.T().Fatalf("Failed to clear words data: %v", err)
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		suite.T().Fatalf("Failed to commit transaction: %v", err)
	}
}

// seedTestData seeds the test database with test data
func (suite *DashboardHandlerTestSuite) seedTestData() {
	// Create test words
	testWords := []models.Word{
		{Portuguese: "olá", English: "hello"},
		{Portuguese: "adeus", English: "goodbye"},
		{Portuguese: "obrigado", English: "thank you"},
	}

	// Insert test words into the database
	for _, word := range testWords {
		stmt, err := suite.db.DB.Prepare("INSERT INTO words (portuguese, english, created_at) VALUES (?, ?, ?)")
		if err != nil {
			suite.T().Fatalf("Failed to prepare statement: %v", err)
		}
		result, err := stmt.Exec(word.Portuguese, word.English, time.Now())
		if err != nil {
			suite.T().Fatalf("Failed to insert test word: %v", err)
		}
		id, _ := result.LastInsertId()
		suite.testWords = append(suite.testWords, &models.Word{
			ID:         id,
			Portuguese: word.Portuguese,
			English:    word.English,
		})
		stmt.Close()
	}

	// Create test groups
	testGroups := []models.Group{
		{Name: "Basics"},
		{Name: "Greetings"},
	}

	// Insert test groups into the database
	for _, group := range testGroups {
		stmt, err := suite.db.DB.Prepare("INSERT INTO groups (name, created_at) VALUES (?, ?)")
		if err != nil {
			suite.T().Fatalf("Failed to prepare statement: %v", err)
		}
		result, err := stmt.Exec(group.Name, time.Now())
		if err != nil {
			suite.T().Fatalf("Failed to insert test group: %v", err)
		}
		id, _ := result.LastInsertId()
		suite.testGroups = append(suite.testGroups, &models.Group{
			ID:   id,
			Name: group.Name,
		})
		stmt.Close()
	}

	// Add words to groups
	if len(suite.testGroups) > 0 && len(suite.testWords) >= 2 {
		for i := 0; i < 2; i++ {
			stmt, err := suite.db.DB.Prepare("INSERT INTO words_groups (word_id, group_id) VALUES (?, ?)")
			if err != nil {
				suite.T().Fatalf("Failed to prepare statement: %v", err)
			}
			_, err = stmt.Exec(suite.testWords[i].ID, suite.testGroups[0].ID)
			if err != nil {
				suite.T().Fatalf("Failed to add word to group: %v", err)
			}
			stmt.Close()
		}
	}

	// Create test study activities
	testActivities := []models.StudyActivity{
		{Name: "Flashcards", ThumbnailURL: "/images/flashcards.png", Description: "Practice with flashcards"},
		{Name: "Quiz", ThumbnailURL: "/images/quiz.png", Description: "Test your knowledge with a quiz"},
	}

	// Insert test activities into the database
	for _, activity := range testActivities {
		stmt, err := suite.db.DB.Prepare("INSERT INTO study_activities (name, thumbnail_url, description, created_at) VALUES (?, ?, ?, ?)")
		if err != nil {
			suite.T().Fatalf("Failed to prepare statement: %v", err)
		}
		result, err := stmt.Exec(activity.Name, activity.ThumbnailURL, activity.Description, time.Now())
		if err != nil {
			suite.T().Fatalf("Failed to insert test activity: %v", err)
		}
		id, _ := result.LastInsertId()
		suite.testStudyActivities = append(suite.testStudyActivities, &models.StudyActivity{
			ID:           id,
			Name:         activity.Name,
			ThumbnailURL: activity.ThumbnailURL,
			Description:  activity.Description,
		})
		stmt.Close()
	}

	// Create test study sessions
	if len(suite.testStudyActivities) > 0 && len(suite.testGroups) > 0 {
		// Create a study session for the first activity and first group
		stmt, err := suite.db.DB.Prepare("INSERT INTO study_sessions (group_id, study_activity_id, created_at) VALUES (?, ?, ?)")
		if err != nil {
			suite.T().Fatalf("Failed to prepare statement: %v", err)
		}

		// Create a session from yesterday
		yesterday := time.Now().AddDate(0, 0, -1)
		result, err := stmt.Exec(suite.testGroups[0].ID, suite.testStudyActivities[0].ID, yesterday)
		if err != nil {
			suite.T().Fatalf("Failed to insert test study session: %v", err)
		}
		id, _ := result.LastInsertId()
		suite.testStudySessions = append(suite.testStudySessions, &models.StudySession{
			ID:              id,
			GroupID:         suite.testGroups[0].ID,
			StudyActivityID: suite.testStudyActivities[0].ID,
			CreatedAt:       yesterday,
		})

		// Create a more recent session (today)
		result, err = stmt.Exec(suite.testGroups[0].ID, suite.testStudyActivities[0].ID, time.Now())
		if err != nil {
			suite.T().Fatalf("Failed to insert test study session: %v", err)
		}
		id, _ = result.LastInsertId()
		suite.testStudySessions = append(suite.testStudySessions, &models.StudySession{
			ID:              id,
			GroupID:         suite.testGroups[0].ID,
			StudyActivityID: suite.testStudyActivities[0].ID,
			CreatedAt:       time.Now(),
		})

		stmt.Close()
	}

	// Create test word review items
	if len(suite.testStudySessions) > 0 && len(suite.testWords) > 0 {
		// Add review items for the first session
		stmt, err := suite.db.DB.Prepare("INSERT INTO word_review_items (study_session_id, word_id, correct, created_at) VALUES (?, ?, ?, ?)")
		if err != nil {
			suite.T().Fatalf("Failed to prepare statement: %v", err)
		}

		// Add some correct and incorrect reviews
		for i, word := range suite.testWords {
			correct := i%2 == 0 // Alternate between correct and incorrect
			result, err := stmt.Exec(suite.testStudySessions[0].ID, word.ID, correct, time.Now())
			if err != nil {
				suite.T().Fatalf("Failed to insert test word review item: %v", err)
			}
			id, _ := result.LastInsertId()
			suite.testWordReviewItems = append(suite.testWordReviewItems, &models.WordReviewItem{
				ID:             id,
				StudySessionID: suite.testStudySessions[0].ID,
				WordID:         word.ID,
				Correct:        correct,
			})
		}

		// Add some reviews for the second session too
		if len(suite.testStudySessions) > 1 {
			for i, word := range suite.testWords {
				if i < 2 { // Only add reviews for the first two words
					correct := i == 0 // First one correct, second one incorrect
					result, err := stmt.Exec(suite.testStudySessions[1].ID, word.ID, correct, time.Now())
					if err != nil {
						suite.T().Fatalf("Failed to insert test word review item: %v", err)
					}
					id, _ := result.LastInsertId()
					suite.testWordReviewItems = append(suite.testWordReviewItems, &models.WordReviewItem{
						ID:             id,
						StudySessionID: suite.testStudySessions[1].ID,
						WordID:         word.ID,
						Correct:        correct,
					})
				}
			}
		}

		stmt.Close()
	}
}

// TestGetLastStudySession tests the GetLastStudySession endpoint
func (suite *DashboardHandlerTestSuite) TestGetLastStudySession() {
	// Perform the request
	w := testutil.PerformRequest(
		suite.T(),
		suite.router,
		"GET",
		"/api/dashboard/last_study_session",
		nil,
	)

	// Print the response body for debugging
	body := w.Body.String()
	suite.T().Logf("Response body: %s", body)

	// Check the status code
	testutil.AssertStatusCode(suite.T(), w, http.StatusOK)

	// Parse the response
	var response models.StudySessionDetail
	testutil.ParseResponse(suite.T(), w, &response)

	// Verify the response - should be the most recent session (second one)
	if len(suite.testStudySessions) > 1 {
		assert.Equal(suite.T(), suite.testStudySessions[1].ID, response.ID)
		assert.Equal(suite.T(), suite.testStudySessions[1].GroupID, response.GroupID)
		assert.Equal(suite.T(), suite.testStudySessions[1].StudyActivityID, response.StudyActivityID)
		assert.NotEmpty(suite.T(), response.ActivityName)
		assert.NotEmpty(suite.T(), response.GroupName)
	}
}

// TestGetStudyProgress tests the GetStudyProgress endpoint
func (suite *DashboardHandlerTestSuite) TestGetStudyProgress() {
	// Perform the request
	w := testutil.PerformRequest(
		suite.T(),
		suite.router,
		"GET",
		"/api/dashboard/study_progress",
		nil,
	)

	// Check the status code
	testutil.AssertStatusCode(suite.T(), w, http.StatusOK)

	// Parse the response
	var response struct {
		TotalWordsStudied   int `json:"total_words_studied"`
		TotalAvailableWords int `json:"total_available_words"`
	}
	testutil.ParseResponse(suite.T(), w, &response)

	// Verify the response
	assert.GreaterOrEqual(suite.T(), response.TotalAvailableWords, len(suite.testWords))
	assert.GreaterOrEqual(suite.T(), response.TotalWordsStudied, 0)
}

// TestGetQuickStats tests the GetQuickStats endpoint
func (suite *DashboardHandlerTestSuite) TestGetQuickStats() {
	// Perform the request
	w := testutil.PerformRequest(
		suite.T(),
		suite.router,
		"GET",
		"/api/dashboard/quick-stats",
		nil,
	)

	// Check the status code
	testutil.AssertStatusCode(suite.T(), w, http.StatusOK)

	// Parse the response
	var response struct {
		SuccessRate        float64 `json:"success_rate"`
		TotalStudySessions int     `json:"total_study_sessions"`
		TotalActiveGroups  int     `json:"total_active_groups"`
		StudyStreakDays    int     `json:"study_streak_days"`
	}
	testutil.ParseResponse(suite.T(), w, &response)

	// Verify the response
	assert.GreaterOrEqual(suite.T(), response.SuccessRate, 0.0)
	assert.LessOrEqual(suite.T(), response.SuccessRate, 100.0)

	// If we have study sessions, we should have some data
	if len(suite.testStudySessions) > 0 {
		assert.Equal(suite.T(), len(suite.testStudySessions), response.TotalStudySessions)
		assert.Greater(suite.T(), response.TotalActiveGroups, 0)
	}
}

// TestGetLastStudySessionEmpty tests the GetLastStudySession endpoint when there are no sessions
func (suite *DashboardHandlerTestSuite) TestGetLastStudySessionEmpty() {
	// Clear all study sessions
	_, err := suite.db.DB.Exec("DELETE FROM study_sessions")
	if err != nil {
		suite.T().Fatalf("Failed to clear study sessions: %v", err)
	}
	suite.testStudySessions = nil

	// Perform the request
	w := testutil.PerformRequest(
		suite.T(),
		suite.router,
		"GET",
		"/api/dashboard/last_study_session",
		nil,
	)

	// Check the status code
	testutil.AssertStatusCode(suite.T(), w, http.StatusOK)

	// Parse the response
	var response map[string]string
	testutil.ParseResponse(suite.T(), w, &response)

	// Verify the response contains the expected message
	assert.Equal(suite.T(), "No study sessions found", response["message"])
}

// TestMain runs the test suite
func TestDashboardHandlerSuite(t *testing.T) {
	suite.Run(t, new(DashboardHandlerTestSuite))
}
