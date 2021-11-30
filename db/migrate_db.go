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
	db.AutoMigrate(&entity.IntroductionMail{})
	db.AutoMigrate(&entity.Report{})
	db.AutoMigrate(&entity.FinancialHelp{})
	db.AutoMigrate(&entity.Death{})
}