package entity

import "time"

type Death struct {
	Id            int       `gorm:"primaryKey;not null;autoIncrement"`
	UserNik       string    `gorm:"not null"`
	Nik           string    `gorm:"not null"`
	Nama          string    `gorm:"not null"`
	Jenis_kelamin bool      `gorm:"not null"`
	Usia          int       `gorm:"not null"`
	Tgl_wafat     time.Time `gorm:"not null"`
	Alamat        string    `gorm:"not null"`
	HistoryId     int       `gorm:"not null"`
	History       History   `gorm:"foreignKey:HistoryId;references:Id" json:"omitempty"`
}

type DeathQuery struct {
	Id            int
	Nik           string
	Nama          string
	Jenis_kelamin bool
	Usia          int
	Tgl_wafat     time.Time
	Alamat        string
	HistoryId     int
	Status        string
}
