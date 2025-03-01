package service

import (
	"github.com/andlogreg/free-genai-bootcamp-2025/backend_go/internal/models"
	"github.com/andlogreg/free-genai-bootcamp-2025/backend_go/internal/repository"
)

type GroupService struct {
	groupRepo *repository.GroupRepository
}

func NewGroupService(groupRepo *repository.GroupRepository) *GroupService {
	return &GroupService{groupRepo: groupRepo}
}

func (s *GroupService) ListGroups() ([]*models.Group, error) {
	return s.groupRepo.ListGroups()
}

// GetGroup returns a group with its stats
func (s *GroupService) GetGroup(id int64) (*models.GroupDetail, error) {
	// Get the basic group info
	group, err := s.groupRepo.GetGroup(id)
	if err != nil {
		return nil, err
	}

	// Get the word count for this group
	wordCount, err := s.groupRepo.CountGroupWords(id)
	if err != nil {
		return nil, err
	}

	// Create a GroupDetail object with the required stats
	return &models.GroupDetail{
		ID:   group.ID,
		Name: group.Name,
		Stats: models.GroupStats{
			TotalWordCount: wordCount,
		},
	}, nil
}

func (s *GroupService) GetGroupWords(id int64) ([]*models.Word, error) {
	return s.groupRepo.GetGroupWords(id)
}

func (s *GroupService) CreateGroup(group *models.Group) (*models.Group, error) {
	return s.groupRepo.CreateGroup(group)
}

func (s *GroupService) UpdateGroup(group *models.Group) (*models.Group, error) {
	return s.groupRepo.UpdateGroup(group)
}

func (s *GroupService) DeleteGroup(id int64) error {
	return s.groupRepo.DeleteGroup(id)
}

func (s *GroupService) AddWordsToGroup(groupID int64, wordIDs []int64) error {
	return s.groupRepo.AddWordsToGroup(groupID, wordIDs)
}

func (s *GroupService) RemoveWordFromGroup(groupID, wordID int64) error {
	return s.groupRepo.RemoveWordFromGroup(groupID, wordID)
}

// ListGroupsPaginated returns a paginated list of groups
func (s *GroupService) ListGroupsPaginated(page, pageSize int) ([]*models.Group, int, error) {
	return s.groupRepo.ListGroupsPaginated(page, pageSize)
}

// GetGroupWordsPaginated returns a paginated list of words in a group with stats
func (s *GroupService) GetGroupWordsPaginated(groupID int64, page, pageSize int) ([]*models.WordWithStats, int, error) {
	return s.groupRepo.GetGroupWordsPaginated(groupID, page, pageSize)
}
