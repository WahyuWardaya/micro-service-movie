package models

import "time"

// Minimal mapping to users table in user-service database
type User struct {
	ID                    uint       `gorm:"primaryKey;column:id" json:"id"`
	Name                  string     `gorm:"column:name" json:"name"`
	Email                 string     `gorm:"column:email" json:"email"`
	SubscriptionType      string     `gorm:"column:subscription_type" json:"subscription_type"`
	SubscriptionExpiredAt *time.Time `gorm:"column:subscription_expired_at" json:"subscription_expired_at"`
	// Note: other columns exist in users table, but we map only the fields we need
}
