package db

import (
	// "fmt"
	// "smartville-server/config"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenDatabase() *gorm.DB {
	// cfg, _ := config.NewConfig(".env")
	// dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
	// 	cfg.Database.Host,
	// 	cfg.Database.Port,
	// 	cfg.Database.Username,
	// 	cfg.Database.Password,
	// 	cfg.Database.Name)
	db, err := gorm.Open(postgres.Open(string(os.Getenv("DATABASE_URL"))), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
