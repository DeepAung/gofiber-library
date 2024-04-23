package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	App App
	DB  DB
}

type App struct {
	Port      string
	JwtSecret string
	GCPBucket string
}

type DB struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string
}

func NewConfig() *Config {
	if len(os.Args) > 1 {
		err := godotenv.Load(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
	}

	return &Config{
		DB: DB{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSL_MODE"),
			TimeZone: os.Getenv("DB_TIME_ZONE"),
		},
		App: App{
			Port:      os.Getenv("APP_PORT"),
			JwtSecret: os.Getenv("APP_JWTSECRET"),
			GCPBucket: os.Getenv("APP_GCPBUCKET"),
		},
	}
}
