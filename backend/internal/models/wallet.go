package models

import "time"

type Wallet struct {
    ID           uint      `gorm:"primaryKey"`
    UserID       uint      `gorm:"uniqueIndex"`
    BalanceCents int64
    UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}
