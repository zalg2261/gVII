package models

import "time"

type Refund struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    BookingID   uint      `json:"booking_id"`
    UserID      uint      `json:"user_id"`
    AmountCents int64     `json:"amount_cents"`
    Reason      string    `json:"reason"`
    Status      string    `json:"status"` // REQUESTED, APPROVED, FAILED
    CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
    ProcessedAt time.Time `json:"processed_at"`
}
