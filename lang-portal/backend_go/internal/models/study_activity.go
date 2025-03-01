package models

import "time"

type StudyActivity struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	ThumbnailURL string    `json:"thumbnail_url"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
}
