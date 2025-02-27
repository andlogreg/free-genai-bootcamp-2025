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
