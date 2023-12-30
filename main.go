package main

import (
	"fmt"
	"log"

	"github.com/DeepAung/gofiber-library/configs"
	"github.com/DeepAung/gofiber-library/modules/books/booksHandler"
	"github.com/DeepAung/gofiber-library/modules/books/booksService"
	"github.com/DeepAung/gofiber-library/modules/users/usersHandler"
	"github.com/DeepAung/gofiber-library/modules/users/usersService"
	"github.com/DeepAung/gofiber-library/modules/views/viewsHandler"
	"github.com/DeepAung/gofiber-library/pkg/databases"
	"github.com/DeepAung/gofiber-library/pkg/middlewares"
	"github.com/DeepAung/gofiber-library/pkg/utils"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/gofiber/template/html/v2"
)

type Server struct {
	App *fiber.App
	Mid *middlewares.Middleware
	Cfg *configs.Config
	DB  *gorm.DB
}

func main() {
	engine := html.New("./views", ".html")

	server := new(Server)
	server.App = fiber.New(fiber.Config{Views: engine})
	server.Mid = middlewares.NewMiddleware()
	server.Cfg = configs.NewConfig()
	server.DB = databases.NewDB(server.Cfg)

	server.App.Use(logger.New())
	server.App.Use(recover.New())
	server.App.Static("/", "./public")

	server.initRoutes()

	addr := fmt.Sprintf("%s:%s", server.Cfg.Fiber.Host, server.Cfg.Fiber.Port)
	log.Fatal(server.App.Listen(addr))
}

func (s *Server) initRoutes() {
	apiGroup := s.App.Group("/api")

	myvalidator := utils.NewMyValidator()
	usersService := usersService.NewUsersService(s.DB, s.Cfg)
	usersHandler.NewUsersHandler(apiGroup, myvalidator, usersService, s.Mid)

	booksService := booksService.NewBooksService(s.DB)
	booksHandler.NewBooksHandler(apiGroup, myvalidator, booksService, usersService, s.Mid)

	viewsHandler.NewViewsHandler(s.App, usersService, booksService, s.Mid, s.Cfg)

	s.App.Use(func(c *fiber.Ctx) error {
		payload, err := usersService.VerifyTokenByCookie(c, "access_token")

		return c.Render("error", fiber.Map{
			"IsAuthenticated": err == nil,
			"Payload":         payload,
		},
			"layouts/main",
		)

	})
}
