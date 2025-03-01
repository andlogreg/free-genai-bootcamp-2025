package models

import "time"

type Word struct {
	ID         int64     `json:"id"`
	Portuguese string    `json:"portuguese"`
	English    string    `json:"english"`
	CreatedAt  time.Time `json:"created_at"`
}

type WordWithStats struct {
	Word
	CorrectCount int `json:"correct_count"`
	WrongCount   int `json:"wrong_count"`
}

// WordStats represents statistics for a word
type WordStats struct {
	CorrectCount int `json:"correct_count"`
	WrongCount   int `json:"wrong_count"`
}

// WordGroup represents a simplified group for inclusion in WordDetail
type WordGroup struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// WordDetail represents a word with its statistics and groups
type WordDetail struct {
	ID         int64       `json:"id"`
	Portuguese string      `json:"portuguese"`
	English    string      `json:"english"`
	Stats      WordStats   `json:"stats"`
	Groups     []WordGroup `json:"groups"`
}
