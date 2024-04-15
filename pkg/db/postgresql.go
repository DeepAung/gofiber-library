package db

import (
	"fmt"
	"log"

	"github.com/DeepAung/gofiber-library/modules/models"
	"github.com/DeepAung/gofiber-library/pkg/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(cfg *configs.Config) *gorm.DB {
	println("connecting to DB")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.PostgreSQL.Host,
		cfg.PostgreSQL.Port,
		cfg.PostgreSQL.User,
		cfg.PostgreSQL.Password,
		cfg.PostgreSQL.DBName,
		cfg.PostgreSQL.SSLMode,
		cfg.PostgreSQL.TimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	println("DB opened")

	err = models.InitModel(db)
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.User{}, &models.Book{})
	println("DB auto migrated")

	return db
}
