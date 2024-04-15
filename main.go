package main

import (
	"github.com/DeepAung/gofiber-library/pkg/configs"
	"github.com/DeepAung/gofiber-library/pkg/db"
	"github.com/DeepAung/gofiber-library/pkg/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{Views: engine})
	mid := middlewares.NewMiddleware()
	cfg := configs.NewConfig()
	db := db.NewDB(cfg)
	server := NewServer(app, mid, cfg, db)

	server.Start()
}
