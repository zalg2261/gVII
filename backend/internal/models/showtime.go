package models

import "time"

type Showtime struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	MovieID    uint      `json:"movie_id"`
	BranchID   uint      `json:"branch_id"`
	Studio     string    `json:"studio"`
	ShowTime   time.Time `json:"show_time"`
	SeatsTotal int       `json:"seats_total" gorm:"default:50"`
	SeatsLeft  int       `json:"seats_left" gorm:"default:50"`
	Price      int64     `json:"price" gorm:"default:50000"` // harga per tiket dalam rupiah
	Status     string    `json:"status" gorm:"default:ACTIVE"` // ACTIVE, CANCELLED
	CreatedAt  time.Time `json:"created_at"`

	Movie  Movie  `json:"movie" gorm:"foreignKey:MovieID"`
	Branch Branch `json:"branch" gorm:"foreignKey:BranchID"`
}
