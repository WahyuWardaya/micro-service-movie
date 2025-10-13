package models

type Genre struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"unique;type:varchar(100)" json:"name"`

	Movies []Movie `gorm:"many2many:movie_genres" json:"-"`
}
