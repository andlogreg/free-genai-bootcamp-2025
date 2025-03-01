package handlers_test

import (
	"encoding/json"
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
	// Make a request to list groups
	w := testutil.PerformRequest(suite.T(), suite.router, "GET", "/api/groups", nil)
	testutil.AssertStatusCode(suite.T(), w, http.StatusOK)

	// Parse the response
	var response models.PaginatedResponse
	testutil.ParseResponse(suite.T(), w, &response)

	// Extract the groups from the items field
	groupsData, err := json.Marshal(response.Items)
	assert.NoError(suite.T(), err)

	var groups []models.Group
	err = json.Unmarshal(groupsData, &groups)
	assert.NoError(suite.T(), err)

	// Verify the response has the expected number of groups
	assert.Len(suite.T(), groups, len(suite.testGroups))

	// Create a map of group IDs to names for easier verification
	groupMap := make(map[int64]string)
	for _, group := range suite.testGroups {
		groupMap[group.ID] = group.Name
	}

	// Verify that our test groups are in the response
	for _, group := range groups {
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
	// Get the first test group
	group := suite.testGroups[0]

	// Make a request to get the group's words
	w := testutil.PerformRequest(suite.T(), suite.router, "GET", fmt.Sprintf("/api/groups/%d/words", group.ID), nil)
	testutil.AssertStatusCode(suite.T(), w, http.StatusOK)

	// Parse the response
	var response models.PaginatedResponse
	testutil.ParseResponse(suite.T(), w, &response)

	// Extract the words from the items field
	wordsData, err := json.Marshal(response.Items)
	assert.NoError(suite.T(), err)

	var words []models.Word
	err = json.Unmarshal(wordsData, &words)
	assert.NoError(suite.T(), err)

	// Verify the response has the expected number of words
	// We expect 2 words for the first group
	assert.Len(suite.T(), words, 2)

	// Create a map of Portuguese words to English translations for easier verification
	wordMap := make(map[string]string)
	for _, word := range suite.testWords {
		wordMap[word.Portuguese] = word.English
	}

	// Verify that our test words are in the response
	for _, word := range words {
		// If this is one of our test words, verify the translation
		if expectedEnglish, ok := wordMap[word.Portuguese]; ok {
			assert.Equal(suite.T(), expectedEnglish, word.English)
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
	// Get the second test group
	group := suite.testGroups[1]

	// Get the third test word
	word := suite.testWords[2]

	// Make a request to add the word to the group
	w := testutil.PerformRequest(
		suite.T(),
		suite.router,
		"POST",
		fmt.Sprintf("/api/groups/%d/words", group.ID),
		[]int64{word.ID},
	)
	testutil.AssertStatusCode(suite.T(), w, http.StatusNoContent)

	// Make a request to get the group's words
	w = testutil.PerformRequest(
		suite.T(),
		suite.router,
		"GET",
		fmt.Sprintf("/api/groups/%d/words", group.ID),
		nil,
	)
	testutil.AssertStatusCode(suite.T(), w, http.StatusOK)

	// Parse the response
	var response models.PaginatedResponse
	testutil.ParseResponse(suite.T(), w, &response)

	// Extract the words from the items field
	wordsData, err := json.Marshal(response.Items)
	assert.NoError(suite.T(), err)

	var words []models.Word
	err = json.Unmarshal(wordsData, &words)
	assert.NoError(suite.T(), err)

	// Verify the response has the expected number of words
	assert.Len(suite.T(), words, 1)

	// Verify that the word is in the response
	assert.Equal(suite.T(), word.ID, words[0].ID)
	assert.Equal(suite.T(), word.Portuguese, words[0].Portuguese)
	assert.Equal(suite.T(), word.English, words[0].English)
}

// TestRemoveWordFromGroup tests the RemoveWordFromGroup endpoint
func (suite *GroupHandlerTestSuite) TestRemoveWordFromGroup() {
	// Get the first test group
	group := suite.testGroups[0]

	// Get the first test word
	word := suite.testWords[0]

	// Make a request to remove the word from the group
	w := testutil.PerformRequest(
		suite.T(),
		suite.router,
		"DELETE",
		fmt.Sprintf("/api/groups/%d/words/%d", group.ID, word.ID),
		nil,
	)
	testutil.AssertStatusCode(suite.T(), w, http.StatusNoContent)

	// Make a request to get the group's words
	w = testutil.PerformRequest(
		suite.T(),
		suite.router,
		"GET",
		fmt.Sprintf("/api/groups/%d/words", group.ID),
		nil,
	)
	testutil.AssertStatusCode(suite.T(), w, http.StatusOK)

	// Parse the response
	var response models.PaginatedResponse
	testutil.ParseResponse(suite.T(), w, &response)

	// Extract the words from the items field
	wordsData, err := json.Marshal(response.Items)
	assert.NoError(suite.T(), err)

	var words []models.Word
	err = json.Unmarshal(wordsData, &words)
	assert.NoError(suite.T(), err)

	// Verify the response has the expected number of words
	assert.Len(suite.T(), words, 1)

	// Verify that the remaining word is the second word
	assert.Equal(suite.T(), suite.testWords[1].ID, words[0].ID)
}

// TestMain runs the test suite
func TestGroupHandlerSuite(t *testing.T) {
	suite.Run(t, new(GroupHandlerTestSuite))
}
