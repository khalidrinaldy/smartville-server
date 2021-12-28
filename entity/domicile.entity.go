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
	HistoryId       int       `gorm:"not null"`
	History         History   `gorm:"foreignKey:HistoryId;references:Id" json:"omitempty"`
}

type DomicileQuery struct {
	Id              int
	Nik_pemohon     string
	Nama_pemohon    string
	Tgl_lahir       time.Time
	Asal_domisili   string
	Tujuan_domisili string
	Status          string
}
