package db

import (
	"fmt"
	"log"

	"github.com/DeepAung/gofiber-library/pkg/configs"
	"github.com/DeepAung/gofiber-library/types"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(cfg *configs.Config) *gorm.DB {
	println("connecting to DB")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.DBName,
		cfg.DB.SSLMode,
		cfg.DB.TimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	println("DB opened")

	err = db.SetupJoinTable(&types.User{}, "FavBooks", &types.UserFavbooks{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&types.User{}, &types.Book{}, &types.Attachment{})
	println("DB auto migrated")

	return db
}
