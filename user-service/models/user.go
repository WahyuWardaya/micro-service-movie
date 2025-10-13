package models

type User struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Name         string `gorm:"type:varchar(100)" json:"name"`
	Email        string `gorm:"uniqueIndex;type:varchar(100)" json:"email"`
	PasswordHash string `gorm:"type:text" json:"-"`
}
