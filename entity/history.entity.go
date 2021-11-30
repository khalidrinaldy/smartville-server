package entity

import "time"

type History struct {
	Id        int       `gorm:"primaryKey;not null;autoIncrement"`
	UserNik   string    `gorm:"not null"`
	Deskripsi string    `gorm:"not null"`
	Waktu     time.Time `gorm:"not null"`
	Status    string    `gorm:"not null"`
}
