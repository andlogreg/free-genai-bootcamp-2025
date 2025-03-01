package models

import "time"

type Group struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type GroupWithStats struct {
	Group
	WordCount int `json:"word_count"`
}

// GroupStats represents statistics for a group
type GroupStats struct {
	TotalWordCount int `json:"total_word_count"`
}

// GroupDetail represents a group with its statistics
type GroupDetail struct {
	ID    int64      `json:"id"`
	Name  string     `json:"name"`
	Stats GroupStats `json:"stats"`
}
