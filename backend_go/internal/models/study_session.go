package models

import "time"

type StudySession struct {
	ID               int64     `json:"id"`
	StudyActivityID  int64     `json:"study_activity_id"`
	GroupID          int64     `json:"group_id"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	CreatedAt        time.Time `json:"created_at"`
	ActivityName     string    `json:"activity_name"`
	GroupName        string    `json:"group_name"`
	ReviewItemsCount int       `json:"review_items_count"`
}

type StudySessionDetail struct {
	ID               int64      `json:"id"`
	ActivityName     string     `json:"activity_name"`
	GroupName        string     `json:"group_name"`
	CreatedAt        time.Time  `json:"start_time"`
	EndTime          *time.Time `json:"end_time,omitempty"`
	ReviewItemsCount int        `json:"review_items_count"`
	StudyActivityID  int64      `json:"-"`
	GroupID          int64      `json:"-"`
}
