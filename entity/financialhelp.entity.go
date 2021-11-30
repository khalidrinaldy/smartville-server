package entity

type FinancialHelp struct {
	Id                int    `gorm:"primaryKey;not null;autoIncrement"`
	UserNik           string `gorm:"not null"`
	Nama_bantuan      string `gorm:"not null"`
	Jenis_bantuan     string `gorm:"not null"`
	Jumlah_dana       int    `gorm:"not null"`
	Alokasi_dana      int    `gorm:"not null"`
	Dana_terealisasi  int    `gorm:"not null"`
	Sisa_dana_bantuan int    `gorm:"not null"`
}
