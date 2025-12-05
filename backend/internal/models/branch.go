package models

import "time"

type Branch struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	City      string    `json:"city"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
}
