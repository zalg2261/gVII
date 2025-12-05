package models

import "time"

type Movie struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title"`
	Genre     string    `json:"genre"`
	Duration  int       `json:"duration"` // dalam menit
	Synopsis  string    `json:"synopsis" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
}
