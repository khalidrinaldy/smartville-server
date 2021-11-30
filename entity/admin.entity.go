package entity

type Admin struct {
	Id          int64  `gorm:"primaryKey;not null;autoIncrement"`
	Nama        string `gorm:"not null"`
	Email       string `gorm:"not null"`
	Password    string `gorm:"not null"`
	Profile_pic string
	Role        string `gorm:"default:admin"`
	Token       string `gorm:"not null"`
}
