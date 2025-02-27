package service

import (
	"github.com/andlogreg/free-genai-bootcamp-2025/backend_go/internal/models"
	"github.com/andlogreg/free-genai-bootcamp-2025/backend_go/internal/repository"
)

type DashboardService struct {
	studySessionRepo *repository.StudySessionRepository
	wordRepo         *repository.WordRepository
	groupRepo        *repository.GroupRepository
}

func NewDashboardService(
	studySessionRepo *repository.StudySessionRepository,
	wordRepo *repository.WordRepository,
	groupRepo *repository.GroupRepository,
) *DashboardService {
	return &DashboardService{
		studySessionRepo: studySessionRepo,
		wordRepo:         wordRepo,
		groupRepo:        groupRepo,
	}
}

func (s *DashboardService) GetLastStudySession() (*models.StudySessionDetail, error) {
	return s.studySessionRepo.GetLastStudySession()
}

type StudyProgress struct {
	TotalWordsStudied   int `json:"total_words_studied"`
	TotalAvailableWords int `json:"total_available_words"`
}

func (s *DashboardService) GetStudyProgress() (*StudyProgress, error) {
	// Get total available words
	totalWords, err := s.wordRepo.CountWords()
	if err != nil {
		return nil, err
	}

	// Get total words studied (distinct words that have been reviewed)
	totalStudied, err := s.studySessionRepo.GetTotalDistinctWordsStudied()
	if err != nil {
		return nil, err
	}

	return &StudyProgress{
		TotalWordsStudied:   totalStudied,
		TotalAvailableWords: totalWords,
	}, nil
}

type QuickStats struct {
	SuccessRate        float64 `json:"success_rate"`
	TotalStudySessions int     `json:"total_study_sessions"`
	TotalActiveGroups  int     `json:"total_active_groups"`
	StudyStreakDays    int     `json:"study_streak_days"`
}

func (s *DashboardService) GetQuickStats() (*QuickStats, error) {
	// Get success rate
	correct, total, err := s.studySessionRepo.GetWordReviewStats()
	if err != nil {
		return nil, err
	}

	var successRate float64
	if total > 0 {
		successRate = float64(correct) * 100 / float64(total)
	}

	// Get total study sessions
	totalSessions, err := s.studySessionRepo.CountStudySessions()
	if err != nil {
		return nil, err
	}

	// Get total active groups
	activeGroups, err := s.studySessionRepo.GetTotalActiveGroups()
	if err != nil {
		return nil, err
	}

	// Get study streak
	streak, err := s.studySessionRepo.GetStudyStreak()
	if err != nil {
		return nil, err
	}

	return &QuickStats{
		SuccessRate:        successRate,
		TotalStudySessions: totalSessions,
		TotalActiveGroups:  activeGroups,
		StudyStreakDays:    streak,
	}, nil
}
