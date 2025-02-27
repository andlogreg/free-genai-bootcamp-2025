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
