package service

import (
	"github.com/andlogreg/free-genai-bootcamp-2025/backend_go/internal/models"
	"github.com/andlogreg/free-genai-bootcamp-2025/backend_go/internal/repository"
)

type StudyActivityService struct {
	activityRepo *repository.StudyActivityRepository
	sessionRepo  *repository.StudySessionRepository
}

func NewStudyActivityService(
	activityRepo *repository.StudyActivityRepository,
	sessionRepo *repository.StudySessionRepository,
) *StudyActivityService {
	return &StudyActivityService{
		activityRepo: activityRepo,
		sessionRepo:  sessionRepo,
	}
}

func (s *StudyActivityService) GetStudyActivity(id int64) (*models.StudyActivity, error) {
	return s.activityRepo.GetStudyActivity(id)
}

func (s *StudyActivityService) ListStudyActivities() ([]models.StudyActivity, error) {
	return s.activityRepo.ListStudyActivities()
}

func (s *StudyActivityService) GetStudyActivitySessions(activityID int64, page, perPage int) ([]models.StudySessionDetail, int, error) {
	offset := (page - 1) * perPage

	sessions, err := s.activityRepo.GetStudyActivitySessions(activityID, offset, perPage)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.activityRepo.CountStudyActivitySessions(activityID)
	if err != nil {
		return nil, 0, err
	}

	return sessions, total, nil
}

func (s *StudyActivityService) CreateStudySession(groupID, activityID int64) (*models.StudySession, error) {
	session := &models.StudySession{
		GroupID:         groupID,
		StudyActivityID: activityID,
	}

	err := s.sessionRepo.CreateStudySession(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *StudyActivityService) ListStudySessions(offset, limit int) ([]models.StudySessionDetail, error) {
	return s.sessionRepo.ListStudySessions(offset, limit)
}

func (s *StudyActivityService) CountStudySessions() (int, error) {
	return s.sessionRepo.CountStudySessions()
}
