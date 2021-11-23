package entity

import "time"

type BirthRegistration struct {
	Id                int    `gorm:"primaryKey;not null;autoIncrement"`
	UserNik           string `gorm:"not null"`
	Nama_bayi         string `gorm:"not null"`
	Jenis_kelamin     bool   `gorm:"not null"`
	Nama_ayah         string `gorm:"not null"`
	Nama_ibu          string `gorm:"not null"`
	Anak_ke           int    `gorm:"not null"`
	Tanggal_kelahiran time.Time
	Alamat_kelahiran  string `gorm:"not null"`
}