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
	Alamat_kelahiran  string  `gorm:"not null"`
	HistoryId         int     `gorm:"not null"`
	History           History `gorm:"foreignKey:HistoryId;references:Id" json:"omitempty"`
}

type BirthQuery struct {
	Id                int
	Nama_bayi         string
	Nama_ayah         string
	Nama_ibu          string
	Anak_ke           int
	Tanggal_kelahiran time.Time
	Alamat_kelahiran  string
	Status            string
}
