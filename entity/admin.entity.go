package entity

type Admin struct {
	Id          int64  `gorm:"primaryKey;not null;autoIncrement"`
	Nama        string `gorm:"not null"`
	Username    string `gorm:"not null"`
	Email       string `gorm:"not null"`
	Password    string `gorm:"not null"`
	Profile_pic string
	Role        string `gorm:"default:admin"`
	Token       string `gorm:"not null"`
}

type Post struct {
	Id       int64  `gorm:"primaryKey;not null;autoIncrement"`
	Title    string `gorm:"not null"`
	Admin_id int64
	Admin    Admin `gorm:"foreignKey:admin_id"`
}

type AdminsPost struct {
	Id    int64
	Nama  string
	Title string
}
