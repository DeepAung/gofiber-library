package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgreSQL PostgreSQL
	Fiber      Fiber
	JwtSecret  string
}

type PostgreSQL struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string
}

type Fiber struct {
	Host string
	Port string
}

func NewConfig() *Config {
	if len(os.Args) > 1 {
		err := godotenv.Load(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
	}

	return &Config{
		PostgreSQL: PostgreSQL{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSL_MODE"),
			TimeZone: os.Getenv("DB_TIME_ZONE"),
		},
		Fiber: Fiber{
			Host: os.Getenv("FIBER_HOST"),
			Port: os.Getenv("FIBER_PORT"),
		},
		JwtSecret: os.Getenv("JWT_SECRET"),
	}
}
