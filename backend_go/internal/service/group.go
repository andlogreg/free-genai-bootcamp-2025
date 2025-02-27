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

func (s *GroupService) GetGroup(id int64) (*models.Group, error) {
	return s.groupRepo.GetGroup(id)
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
