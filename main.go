package main

import (
	"log"

	"github.com/DeepAung/gofiber-library/configs"
	"github.com/DeepAung/gofiber-library/modules/books"
	"github.com/DeepAung/gofiber-library/modules/users"
	"github.com/DeepAung/gofiber-library/modules/views"
	"github.com/DeepAung/gofiber-library/pkg/databases"
	"github.com/DeepAung/gofiber-library/pkg/utils"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/gofiber/template/html/v2"
)

type Server struct {
	App *fiber.App
	Cfg *configs.Config
	DB  *gorm.DB
}

func main() {
	engine := html.New("./views", ".html")

	server := new(Server)
	server.App = fiber.New(fiber.Config{Views: engine})
	server.Cfg = configs.NewConfig()
	server.DB = databases.NewDB(server.Cfg)

	server.App.Use(logger.New())
	server.App.Use(recover.New())
	server.App.Static("/", "./public")

	server.initRoutes()

	log.Fatal(server.App.Listen("127.0.0.1:8080"))
}

func (s *Server) initRoutes() {
	apiGroup := s.App.Group("/api")

	myvalidator := utils.NewMyValidator()
	usersService := users.NewUsersService(s.DB, s.Cfg)
	users.NewUsersHandler(apiGroup, myvalidator, usersService)

	booksService := books.NewBooksService(s.DB)
	// books.NewBooksHandler(booksService)

	views.NewViewsHandler(s.App, usersService, booksService, s.Cfg)
}
