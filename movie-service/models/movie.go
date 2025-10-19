package models

import "time"

type Movie struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Title          string    `gorm:"type:varchar(255)" json:"title"`
	PosterBase64   string    `gorm:"type:text" json:"poster_base64"`
	DurationMinutes int       `json:"duration_minutes"`
	Synopsis       string    `gorm:"type:text" json:"synopsis"`
	ReleaseYear    int       `json:"release_year"`
	Rating         float32   `json:"rating"`
	Views          int64     `json:"views"`
	IsPremium      bool      `gorm:"default:true" json:"is_premium"`


	Genres []Genre `gorm:"many2many:movie_genres" json:"genres"`
	Actors []Actor `gorm:"many2many:movie_actors" json:"actors"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
