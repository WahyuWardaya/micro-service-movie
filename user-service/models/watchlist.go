package models

import "time"

// Watchlist akan menjadi join table antara User dan Movie (ID)
type Watchlist struct {
    UserID    uint      `gorm:"primaryKey"`
    MovieID   uint      `gorm:"primaryKey"`
    CreatedAt time.Time
}