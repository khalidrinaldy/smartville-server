package db

import (
	"fmt"
	"os"
	//"smartville-server/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenDatabase() *gorm.DB {
	//cfg, _ := config.NewConfig(".env")
	// dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
	// 	cfg.Database.Host,
	// 	cfg.Database.Port,
	// 	cfg.Database.Username,
	// 	cfg.Database.Password,
	// 	cfg.Database.Name)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_PORT"),
		os.Getenv("SSL_MODE"),
		os.Getenv("TIMEZONE"))
	db, err := gorm.Open(postgres.Open(string(dsn)), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
