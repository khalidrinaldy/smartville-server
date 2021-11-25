package entity

import "time"

type DomicileRegistration struct {
	Id              int       `gorm:"primaryKey;not null;autoIncrement"`
	UserNik         string    `gorm:"not null"`
	Nik_pemohon     string    `gorm:"not null"`
	Nama_pemohon    string    `gorm:"not null"`
	Tgl_lahir       time.Time `gorm:"not null"`
	Asal_domisili   string    `gorm:"not null"`
	Tujuan_domisili string    `gorm:"not null"`
}
