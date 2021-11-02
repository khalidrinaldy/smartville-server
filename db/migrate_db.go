package db

import (
	"smartville-server/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&entity.Admin{})
}