package entity

type IntroductionMail struct {
	Id             int     `gorm:"primaryKey;not null;autoIncrement"`
	UserNik        string  `gorm:"not null"`
	Nik_pemohon    string  `gorm:"not null"`
	Nama_pemohon   string  `gorm:"not null"`
	Alamat_pemohon string  `gorm:"not null"`
	No_hp          string  `gorm:"not null"`
	Jenis_surat    string  `gorm:"not null"`
	HistoryId      int     `gorm:"not null"`
	History        History `gorm:"foreignKey:HistoryId;references:Id" json:"omitempty"`
}

type IntroductionMailQuery struct {
	Id             int
	Nik_pemohon    string
	Nama_pemohon   string
	Alamat_pemohon string
	No_hp          string
	Jenis_surat    string
	Status         string
}
