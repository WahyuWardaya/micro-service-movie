package models

import "time"

type Subscription struct {
    ID uint `gorm:"primaryKey" json:"id"`
    UserID uint `json:"user_id"`
    Plan string `json:"plan"`
    Amount int64 `json:"amount"`
    Status string `json:"status"`
    StartedAt time.Time `json:"start_at"`
    EndAt time.Time `json:"end_at"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
