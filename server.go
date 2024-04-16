package main

import (
	"log"

	"github.com/DeepAung/gofiber-library/handlers"
	"github.com/DeepAung/gofiber-library/pkg/configs"
	"github.com/DeepAung/gofiber-library/pkg/middlewares"
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
	usersSvc := services.NewUsersService(s.DB, s.Cfg)
	booksSvc := services.NewBooksService(s.DB)

	api := s.App.Group("/api")
	s.UsersRouter(api, usersSvc)
	s.BooksRouter(api, booksSvc, usersSvc)
	s.ViewsRouter(s.App, usersSvc, booksSvc)

	s.App.Use(s.Mid.PageNotFound(usersSvc))
}

func (s *server) UsersRouter(r fiber.Router, usersSvc *services.UsersService) {
	handler := handlers.NewUsersHandler(usersSvc, s.Mid)

	onlyAuthorized := s.Mid.JwtAccessTokenAuth(usersSvc)
	onlyUnauthorized := s.Mid.OnlyUnauthorizedAuth(usersSvc)

	r.Post("/login", onlyUnauthorized, handler.Login)
	r.Post("/register", onlyUnauthorized, handler.Register)

	r.Post("/logout", onlyAuthorized, handler.Logout)
	r.Post("/refresh", onlyAuthorized, handler.UpdateTokens)
}

func (s *server) BooksRouter(
	r fiber.Router,
	booksSvc *services.BooksService,
	usersSvc *services.UsersService,
) {
	handler := handlers.NewBooksHandler(booksSvc, s.Mid)

	onlyAuthorized := s.Mid.JwtAccessTokenAuth(usersSvc)
	onlyAdmin := s.Mid.OnlyAdmin(usersSvc)

	book := r.Group("/books", onlyAuthorized)

	book.Post("/", onlyAdmin, handler.CreateBook)
	book.Put("/:id", onlyAdmin, handler.UpdateBook)
	book.Delete("/:id", onlyAdmin, handler.DeleteBook)
	book.Post("/:id/favorite", handler.ToggleFavoriteBook)
}

func (s *server) ViewsRouter(
	r fiber.Router,
	usersSvc *services.UsersService,
	booksSvc *services.BooksService,
) {
	handler := handlers.NewViewsHandler(booksSvc, s.Mid, s.Cfg)

	onlyAuthorized := s.Mid.JwtAccessTokenAuth(usersSvc)
	onlyUnauthorized := s.Mid.OnlyUnauthorizedAuth(usersSvc)

	onlyAdmin := s.Mid.OnlyAdmin(usersSvc)
	setIsAdmin := s.Mid.SetIsAdmin(usersSvc)
	setOnAdminPage := func(c *fiber.Ctx) error {
		c.Locals("onAdminPage", true)
		return c.Next()
	}

	r.Get("/login", onlyUnauthorized, handler.LoginView)
	r.Get("/register", onlyUnauthorized, handler.RegisterView)

	r.Get("/", onlyAuthorized, setIsAdmin, handler.IndexView)
	r.Get("/books/:id", onlyAuthorized, setIsAdmin, handler.DetailView)

	admin := r.Group("/admin", onlyAuthorized, onlyAdmin, setOnAdminPage)
	admin.Get("/", handler.IndexView)
	admin.Get("/books/:id", handler.DetailView)
	admin.Get("/create", handler.CreateView)
}
