package entity

import "time"

type News struct {
	Id               int       `gorm:"primaryKey;not null;autoIncrement"`
	Judul_berita     string    `gorm:"not null"`
	Foto_berita      string    `gorm:"not null"`
	Deskripsi_berita string    `gorm:"not null"`
	Tanggal_terbit   time.Time 
}