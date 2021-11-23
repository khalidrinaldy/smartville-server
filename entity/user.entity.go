package entity

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// type User struct {
// 	Id            int64  `gorm:"primaryKey;not null;autoIncrement"`
// 	Nik           string `gorm:"type:varchar(16);primaryKey;not null"`
// 	Nama          string `gorm:"type:varchar(100);not null"`
// 	Email         string `gorm:"type:varchar(100);not null"`
// 	Password      string `gorm:"type:varchar(255);not null"`
// 	Tgl_lahir     time.Time
// 	Tempat_lahir  string `gorm:"type:varchar(50);not null"`
// 	Alamat        string `gorm:"type:varchar(255);not null"`
// 	Dusun         string `gorm:"type:varchar(50);not null"`
// 	Rt            int    `gorm:"type:int(2);not null"`
// 	Rw            int    `gorm:"type:int(2);not null"`
// 	Jenis_kelamin bool   `gorm:"not null"`
// 	No_hp         string `gorm:"type:varchar(15);not null"`
// 	Role          string `gorm:"type:varchar(15);default:user"`
// 	Profile_pic   string `gorm:"type:varchar(300);not null"`
// 	Token         string `gorm:"type:varchar(255);not null"`
// }

type User struct {
	Id                 int64  `gorm:"primaryKey;not null;autoIncrement"`
	Nik                string `gorm:"primaryKey;not null"`
	Nama               string `gorm:"not null"`
	Email              string `gorm:"not null"`
	Password           string `gorm:"not null"`
	Tgl_lahir          time.Time
	Tempat_lahir       string              `gorm:"not null"`
	Alamat             string              `gorm:"not null"`
	Dusun              string              `gorm:"not null"`
	Rt                 int                 `gorm:"not null"`
	Rw                 int                 `gorm:"not null"`
	Jenis_kelamin      bool                `gorm:"not null"`
	No_hp              string              `gorm:"not null"`
	Role               string              `gorm:"default:user"`
	Profile_pic        string              `gorm:"not null"`
	Token              string              `gorm:"not null"`
	BirthRegistrations []BirthRegistration `gorm:"foreignKey:UserNik;references:Nik" json:",omitempty"`
}

type Claims struct {
	Nik string `json:"username"`
	jwt.StandardClaims
}

type UserList struct {
	Nik  string
	Nama string
}
