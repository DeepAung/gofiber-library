package main

import (
	"log"

	"github.com/DeepAung/gofiber-library/handlers"
	"github.com/DeepAung/gofiber-library/pkg/configs"
	"github.com/DeepAung/gofiber-library/pkg/middlewares"
	"github.com/DeepAung/gofiber-library/pkg/utils"
	"github.com/DeepAung/gofiber-library/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

type server struct {
	App *fiber.App
	Mid *middlewares.Middleware
	Cfg *configs.Config
	DB  *gorm.DB
}

func NewServer(
	app *fiber.App,
	mid *middlewares.Middleware,
	cfg *configs.Config,
	db *gorm.DB,
) *server {
	return &server{
		App: app,
		Mid: mid,
		Cfg: cfg,
		DB:  db,
	}
}

func (s *server) Start() {
	s.App.Use(logger.New())
	s.App.Use(recover.New())
	s.App.Use(cors.New(cors.Config{}))
	s.App.Static("/static", "./static")

	s.initRoutes()

	log.Fatal(s.App.Listen(s.Cfg.Fiber.Port))
}

func (s *server) initRoutes() {
	api := s.App.Group("/api")

	myvalidator := utils.NewMyValidator()
	myerror := utils.NewMyError()

	usersService := services.NewUsersService(s.DB, s.Cfg)
	handlers.NewUsersHandler(api, myvalidator, myerror, usersService, s.Mid)

	booksService := services.NewBooksService(s.DB)
	handlers.NewBooksHandler(api, myvalidator, myerror, booksService, usersService, s.Mid)

	handlers.NewViewsHandler(s.App, myerror, usersService, booksService, s.Mid, s.Cfg)

	s.App.Use(s.Mid.PageNotFound(usersService))
}
