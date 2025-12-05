package models

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"uniqueIndex"`
	Password  string    `json:"-"`     // jangan tampilkan password
	Role      string    `json:"role" gorm:"default:user"` // 'user' or 'admin'
	Balance   int64     `json:"balance" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
}
