package handlers_test

import (
	"fmt"
	"net/http"
	"testing"

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

// GroupHandlerTestSuite is a test suite for the group handlers
type GroupHandlerTestSuite struct {
	suite.Suite
	router       *gin.Engine
	db           *database.TestDB
	groupHandler *handlers.GroupHandler
	testGroups   []*models.Group
	testWords    []*models.Word
}

// SetupSuite sets up the test suite
func (suite *GroupHandlerTestSuite) SetupSuite() {
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
	suite.groupHandler = handlers.NewGroupHandler(groupService)
	studyActivityHandler := handlers.NewStudyActivityHandler(studyActivityService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)

	// Setup router
	suite.router = api.SetupRouter(
		dashboardHandler,
		studyActivityHandler,
		wordHandler,
		suite.groupHandler,
	)
}

// TearDownSuite tears down the test suite
func (suite *GroupHandlerTestSuite) TearDownSuite() {
	// Close and remove the test database
	if suite.db != nil {
		suite.db.Close()
	}
}

// SetupTest sets up each test
func (suite *GroupHandlerTestSuite) SetupTest() {
	// Clear any existing test data
	suite.clearTestData()

	// Seed test data for this test
	suite.seedTestData()
}

// TearDownTest cleans up after each test
func (suite *GroupHandlerTestSuite) TearDownTest() {
	// Clear test data
	suite.clearTestData()
}

// clearTestData removes all test data from the database
func (suite *GroupHandlerTestSuite) clearTestData() {
	// Clear the testGroups and testWords slices
	suite.testGroups = nil
	suite.testWords = nil

	// Use a transaction to handle errors gracefully
	tx, err := suite.db.DB.Begin()
	if err != nil {
		suite.T().Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Check if the words_groups table exists
	var tableExists int
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
func (suite *GroupHandlerTestSuite) seedTestData() {
	// Create test words
	testWords := []models.Word{
		{Portuguese: "olá", English: "hello"},
		{Portuguese: "adeus", English: "goodbye"},
		{Portuguese: "obrigado", English: "thank you"},
	}

	// Insert test words into the database
	for _, word := range testWords {
		stmt, err := suite.db.DB.Prepare("INSERT INTO words (portuguese, english) VALUES (?, ?)")
		if err != nil {
			suite.T().Fatalf("Failed to prepare statement: %v", err)
		}
		result, err := stmt.Exec(word.Portuguese, word.English)
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
		{Name: "Advanced"},
	}

	// Insert test groups into the database
	for _, group := range testGroups {
		stmt, err := suite.db.DB.Prepare("INSERT INTO groups (name) VALUES (?)")
		if err != nil {
			suite.T().Fatalf("Failed to prepare statement: %v", err)
		}
		result, err := stmt.Exec(group.Name)
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
	// Add first two words to the first group
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
}

// TestListGroups tests the ListGroups endpoint
func (suite *GroupHandlerTestSuite) TestListGroups() {
	// Perform the request
	w := testutil.PerformRequest(suite.T(), suite.router, "GET", "/api/groups", nil)

	// Check the status code
	testutil.AssertStatusCode(suite.T(), w, http.StatusOK)

	// Parse the response
	var response []models.Group
	testutil.ParseResponse(suite.T(), w, &response)

	// Verify the response has the expected number of groups
	assert.Len(suite.T(), response, len(suite.testGroups))

	// Create a map of group IDs to names for easier verification
	groupMap := make(map[int64]string)
	for _, group := range suite.testGroups {
		groupMap[group.ID] = group.Name
	}

	// Verify that our test groups are in the response
	for _, group := range response {
		// If this is one of our test groups, verify the name
		if expectedName, ok := groupMap[group.ID]; ok {
			assert.Equal(suite.T(), expectedName, group.Name)
		}
	}
}

// TestGetGroup tests the GetGroup endpoint
func (suite *GroupHandlerTestSuite) TestGetGroup() {
	// Perform the request
	w := testutil.PerformRequest(
		suite.T(),
		suite.router,
		"GET",
		fmt.Sprintf("/api/groups/%d", suite.testGroups[0].ID),
		nil,
	)

	// Check the status code
	testutil.AssertStatusCode(suite.T(), w, http.StatusOK)

	// Parse the response
	var response models.Group
	testutil.ParseResponse(suite.T(), w, &response)

	// Verify the response
	assert.Equal(suite.T(), suite.testGroups[0].ID, response.ID)
	assert.Equal(suite.T(), suite.testGroups[0].Name, response.Name)
}

// TestGetGroupNotFound tests the GetGroup endpoint with a non-existent ID
func (suite *GroupHandlerTestSuite) TestGetGroupNotFound() {
	// Perform the request
	w := testutil.PerformRequest(suite.T(), suite.router, "GET", "/api/groups/9999", nil)

	// Check the status code - should be 404 for non-existent group
	testutil.AssertStatusCode(suite.T(), w, http.StatusNotFound)
}

// TestGetGroupWords tests the GetGroupWords endpoint
func (suite *GroupHandlerTestSuite) TestGetGroupWords() {
	// Perform the request
	w := testutil.PerformRequest(
		suite.T(),
		suite.router,
		"GET",
		fmt.Sprintf("/api/groups/%d/words", suite.testGroups[0].ID),
		nil,
	)

	// Check the status code
	testutil.AssertStatusCode(suite.T(), w, http.StatusOK)

	// Parse the response
	var response []models.Word
	testutil.ParseResponse(suite.T(), w, &response)

	// Verify the response has the expected number of words (first two words)
	assert.Len(suite.T(), response, 2)

	// Create a map of word IDs to Portuguese words for easier verification
	wordMap := make(map[int64]string)
	for i := 0; i < 2; i++ {
		wordMap[suite.testWords[i].ID] = suite.testWords[i].Portuguese
	}

	// Verify that the expected words are in the response
	for _, word := range response {
		// If this is one of our test words, verify the Portuguese
		if expectedPortuguese, ok := wordMap[word.ID]; ok {
			assert.Equal(suite.T(), expectedPortuguese, word.Portuguese)
		}
	}
}

// TestCreateGroup tests the CreateGroup endpoint
func (suite *GroupHandlerTestSuite) TestCreateGroup() {
	// Create a new group
	newGroup := models.Group{
		Name: "New Test Group",
	}

	// Perform the request
	w := testutil.PerformRequest(suite.T(), suite.router, "POST", "/api/groups", newGroup)

	// Check the status code
	testutil.AssertStatusCode(suite.T(), w, http.StatusCreated)

	// Parse the response
	var response models.Group
	testutil.ParseResponse(suite.T(), w, &response)

	// Verify the response
	assert.NotZero(suite.T(), response.ID)
	assert.Equal(suite.T(), newGroup.Name, response.Name)
	assert.NotZero(suite.T(), response.CreatedAt)

	// Verify the group was created in the database
	var count int
	err := suite.db.DB.QueryRow("SELECT COUNT(*) FROM groups WHERE name = ?", newGroup.Name).Scan(&count)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, count)
}

// TestUpdateGroup tests the UpdateGroup endpoint
func (suite *GroupHandlerTestSuite) TestUpdateGroup() {
	// Update the first test group
	updatedGroup := models.Group{
		Name: "Updated Group Name",
	}

	// Perform the request
	w := testutil.PerformRequest(
		suite.T(),
		suite.router,
		"PUT",
		fmt.Sprintf("/api/groups/%d", suite.testGroups[0].ID),
		updatedGroup,
	)

	// Check the status code
	testutil.AssertStatusCode(suite.T(), w, http.StatusOK)

	// Parse the response
	var response models.Group
	testutil.ParseResponse(suite.T(), w, &response)

	// Verify the response
	assert.Equal(suite.T(), suite.testGroups[0].ID, response.ID)
	assert.Equal(suite.T(), updatedGroup.Name, response.Name)

	// Verify the group was updated in the database
	var name string
	err := suite.db.DB.QueryRow("SELECT name FROM groups WHERE id = ?", suite.testGroups[0].ID).Scan(&name)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), updatedGroup.Name, name)
}

// TestDeleteGroup tests the DeleteGroup endpoint
func (suite *GroupHandlerTestSuite) TestDeleteGroup() {
	// Perform the request
	w := testutil.PerformRequest(
		suite.T(),
		suite.router,
		"DELETE",
		fmt.Sprintf("/api/groups/%d", suite.testGroups[2].ID),
		nil,
	)

	// Check the status code
	testutil.AssertStatusCode(suite.T(), w, http.StatusNoContent)

	// Verify the group was deleted from the database
	var count int
	err := suite.db.DB.QueryRow("SELECT COUNT(*) FROM groups WHERE id = ?", suite.testGroups[2].ID).Scan(&count)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, count)
}

// TestAddWordsToGroup tests the AddWordsToGroup endpoint
func (suite *GroupHandlerTestSuite) TestAddWordsToGroup() {
	// Add the third word to the second group
	wordIDs := []int64{suite.testWords[2].ID}

	// Perform the request
	w := testutil.PerformRequest(
		suite.T(),
		suite.router,
		"POST",
		fmt.Sprintf("/api/groups/%d/words", suite.testGroups[1].ID),
		wordIDs,
	)

	// Check the status code
	testutil.AssertStatusCode(suite.T(), w, http.StatusNoContent)

	// Verify the word was added to the group in the database
	var count int
	err := suite.db.DB.QueryRow(
		"SELECT COUNT(*) FROM words_groups WHERE group_id = ? AND word_id = ?",
		suite.testGroups[1].ID, suite.testWords[2].ID,
	).Scan(&count)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, count)

	// Now check that the word appears in the group's words
	w = testutil.PerformRequest(
		suite.T(),
		suite.router,
		"GET",
		fmt.Sprintf("/api/groups/%d/words", suite.testGroups[1].ID),
		nil,
	)
	testutil.AssertStatusCode(suite.T(), w, http.StatusOK)

	var response []models.Word
	testutil.ParseResponse(suite.T(), w, &response)

	// Should have one word
	assert.Len(suite.T(), response, 1)

	// And it should be the third word
	assert.Equal(suite.T(), suite.testWords[2].ID, response[0].ID)
}

// TestRemoveWordFromGroup tests the RemoveWordFromGroup endpoint
func (suite *GroupHandlerTestSuite) TestRemoveWordFromGroup() {
	// First, ensure we have a word in a group
	// We'll use the first word in the first group, which was added in seedTestData

	// Perform the request to remove the word
	w := testutil.PerformRequest(
		suite.T(),
		suite.router,
		"DELETE",
		fmt.Sprintf("/api/groups/%d/words/%d", suite.testGroups[0].ID, suite.testWords[0].ID),
		nil,
	)

	// Check the status code
	testutil.AssertStatusCode(suite.T(), w, http.StatusNoContent)

	// Verify the word was removed from the group in the database
	var count int
	err := suite.db.DB.QueryRow(
		"SELECT COUNT(*) FROM words_groups WHERE group_id = ? AND word_id = ?",
		suite.testGroups[0].ID, suite.testWords[0].ID,
	).Scan(&count)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, count)

	// Now check that the word no longer appears in the group's words
	w = testutil.PerformRequest(
		suite.T(),
		suite.router,
		"GET",
		fmt.Sprintf("/api/groups/%d/words", suite.testGroups[0].ID),
		nil,
	)
	testutil.AssertStatusCode(suite.T(), w, http.StatusOK)

	var response []models.Word
	testutil.ParseResponse(suite.T(), w, &response)

	// Should have one word left (the second word)
	assert.Len(suite.T(), response, 1)

	// And it should be the second word
	assert.Equal(suite.T(), suite.testWords[1].ID, response[0].ID)
}

// TestMain runs the test suite
func TestGroupHandlerSuite(t *testing.T) {
	suite.Run(t, new(GroupHandlerTestSuite))
}
