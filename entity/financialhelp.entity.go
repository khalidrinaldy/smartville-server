package entity

type FinancialHelp struct {
	Id                int     `gorm:"primaryKey;not null;autoIncrement"`
	UserNik           string  `gorm:"not null"`
	Nama_bantuan      string  `gorm:"not null"`
	Jenis_bantuan     string  `gorm:"not null"`
	Jumlah_dana       int     `gorm:"not null"`
	Alokasi_dana      int     `gorm:"not null"`
	Dana_terealisasi  int     `gorm:"not null"`
	Sisa_dana_bantuan int     `gorm:"not null"`
	HistoryId         int     `gorm:"not null"`
	History           History `gorm:"foreignKey:HistoryId;references:Id" json:"omitempty"`
}

type FinancialQuery struct {
	Id                int
	Nama_bantuan      string
	Jenis_bantuan     string
	Jumlah_dana       int
	Alokasi_dana      int
	Dana_terealisasi  int
	Sisa_dana_bantuan int
	HistoryId         int
	Status            string
}
