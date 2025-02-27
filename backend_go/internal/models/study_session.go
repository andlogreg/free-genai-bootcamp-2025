package models

import "time"

type StudySession struct {
	ID              int64     `json:"id"`
	GroupID         int64     `json:"group_id"`
	StudyActivityID int64     `json:"study_activity_id"`
	CreatedAt       time.Time `json:"created_at"`
}

type StudySessionDetail struct {
	StudySession
	ActivityName     string     `json:"activity_name"`
	GroupName        string     `json:"group_name"`
	ReviewItemsCount int        `json:"review_items_count"`
	EndTime          *time.Time `json:"end_time,omitempty"`
}
