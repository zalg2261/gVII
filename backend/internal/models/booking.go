package models

import "time"

type Booking struct {
	ID         uint       `json:"id" gorm:"primaryKey"`
	UserID     uint       `json:"user_id"`
	ShowtimeID uint       `json:"showtime_id"`
	Seats      int        `json:"seats"`
	Status     string     `json:"status"`
	ExpiresAt  *time.Time `json:"expires_at"`
	CreatedAt  time.Time  `json:"created_at"`

	Showtime Showtime `json:"showtime" gorm:"foreignKey:ShowtimeID"`
	User     User     `json:"user" gorm:"foreignKey:UserID"`
}
