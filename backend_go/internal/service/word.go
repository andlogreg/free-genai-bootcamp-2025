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

func (s *WordService) GetWord(id int64) (*models.Word, error) {
	return s.wordRepo.GetWord(id)
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
