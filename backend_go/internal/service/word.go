package service

import (
	"github.com/andlogreg/free-genai-bootcamp-2025/backend_go/internal/models"
	"github.com/andlogreg/free-genai-bootcamp-2025/backend_go/internal/repository"
)

type WordService struct {
	wordRepo *repository.WordRepository
}

func NewWordService(wordRepo *repository.WordRepository) *WordService {
	return &WordService{wordRepo: wordRepo}
}

func (s *WordService) ListWords() ([]*models.Word, error) {
	return s.wordRepo.ListWords()
}

// ListWordsWithStatsPaginated returns a paginated list of words with their stats
func (s *WordService) ListWordsWithStatsPaginated(page, pageSize int) ([]*models.WordWithStats, int, error) {
	return s.wordRepo.ListWordsWithStatsPaginated(page, pageSize)
}

func (s *WordService) GetWord(id int64) (*models.Word, error) {
	return s.wordRepo.GetWord(id)
}

func (s *WordService) GetWordWithStats(id int64) (*models.WordWithStats, error) {
	return s.wordRepo.GetWordWithStats(id)
}

func (s *WordService) CreateWord(word *models.Word) (*models.Word, error) {
	return s.wordRepo.CreateWord(word)
}

func (s *WordService) UpdateWord(word *models.Word) (*models.Word, error) {
	return s.wordRepo.UpdateWord(word)
}

func (s *WordService) DeleteWord(id int64) error {
	return s.wordRepo.DeleteWord(id)
}

// GetWordDetail returns a word with its statistics and groups
func (s *WordService) GetWordDetail(id int64) (*models.WordDetail, error) {
	// Get the word with stats
	wordWithStats, err := s.wordRepo.GetWordWithStats(id)
	if err != nil {
		return nil, err
	}
	if wordWithStats == nil {
		return nil, nil
	}

	// Get the groups for this word
	groups, err := s.wordRepo.GetWordGroups(id)
	if err != nil {
		return nil, err
	}

	// Create a WordDetail object
	wordDetail := &models.WordDetail{
		ID:         wordWithStats.ID,
		Portuguese: wordWithStats.Portuguese,
		English:    wordWithStats.English,
		Stats: models.WordStats{
			CorrectCount: wordWithStats.CorrectCount,
			WrongCount:   wordWithStats.WrongCount,
		},
		Groups: groups,
	}

	return wordDetail, nil
}
