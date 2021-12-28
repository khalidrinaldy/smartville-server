package entity

import "time"

type Report struct {
	Id            int    `gorm:"primaryKey;not null;autoIncrement"`
	UserNik       string `gorm:"not null"`
	Nama_pelapor  string `gorm:"not null"`
	Deskripsi     string `gorm:"not null"`
	Jenis_laporan string `gorm:"not null"`
	Tgl_laporan   time.Time
	No_hp         string  `gorm:"not null"`
	Alamat        string  `gorm:"not null"`
	Foto_kejadian string  `gorm:"not null"`
	HistoryId     int     `gorm:"not null"`
	History       History `gorm:"foreignKey:HistoryId;references:Id" json:"omitempty"`
}

type ReportQuery struct {
	Id            int
	Nama_pelapor  string
	Deskripsi     string
	Jenis_laporan string
	Tgl_laporan   time.Time
	No_hp         string
	Alamat        string
	Foto_kejadian string
	HistoryId     int
	Status        string
}
