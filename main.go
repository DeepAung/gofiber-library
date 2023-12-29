package main

import (
	"fmt"
	"log"

	"github.com/DeepAung/gofiber-library/configs"
	"github.com/DeepAung/gofiber-library/modules/books"
	"github.com/DeepAung/gofiber-library/modules/dummy"
	"github.com/DeepAung/gofiber-library/modules/users"
	"github.com/DeepAung/gofiber-library/modules/views"
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
	usersService := users.NewUsersService(s.DB, s.Cfg)
	users.NewUsersHandler(apiGroup, myvalidator, usersService)

	booksService := books.NewBooksService(s.DB)
	books.NewBooksHandler(apiGroup, myvalidator, booksService, usersService)

	dummyGroup := s.App.Group("/dummy")
	dummy.NewDummyHandler(dummyGroup, s.DB)

	views.NewViewsHandler(s.App, usersService, booksService, s.Mid, s.Cfg)

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
