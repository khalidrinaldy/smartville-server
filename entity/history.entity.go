package entity

import "time"

type History struct {
	Id                 int       `gorm:"primaryKey;not null;autoIncrement"`
	UserNik            string    `gorm:"not null"`
	Perihal            string    `gorm:"not null"`
	Deskripsi          string    `gorm:"not null"`
	CreatedAt          time.Time `gorm:"not null"`
	UpdatedAt          time.Time `gorm:"not null"`
	Status             string    `gorm:"not null"`
	Registration_token string    `gorm:"not null"`
}
