package db

import (
	"smartville-server/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&entity.Admin{})
	db.AutoMigrate(&entity.User{})
	db.Exec("alter table users add constraint users_unique unique (nik)")
	db.AutoMigrate(&entity.News{})
	db.AutoMigrate(&entity.BirthRegistration{})
	db.AutoMigrate(&entity.DomicileRegistration{})
}