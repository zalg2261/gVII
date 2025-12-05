package models

import "time"

type SeatLock struct {
	ID         uint       `json:"id" gorm:"primaryKey"`
	ShowtimeID uint       `json:"showtime_id"`
	BookingID  *uint      `json:"booking_id"`
	SeatsCount int        `json:"seats_count" gorm:"default:1"`
	LockedAt   time.Time  `json:"locked_at"`
	ExpiresAt  time.Time  `json:"expires_at"`
}
