package models

type Actor struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"type:varchar(100)" json:"name"`
	PhotoBase64  string `gorm:"type:text" json:"photo_url"`

	Movies []Movie `gorm:"many2many:movie_actors" json:"-"`
}
