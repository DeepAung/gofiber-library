package views

import (
	"fmt"

	"github.com/DeepAung/gofiber-library/configs"
	"github.com/DeepAung/gofiber-library/modules/books"
	"github.com/DeepAung/gofiber-library/modules/users"
	"github.com/DeepAung/gofiber-library/pkg/middlewares"
	"github.com/gofiber/fiber/v2"
)

type ViewsHandler struct {
	usersService *users.UsersService
	booksService *books.BooksService
	cfg          *configs.Config
}

func NewViewsHandler(
	r fiber.Router,
	usersService *users.UsersService,
	booksService *books.BooksService,
	cfg *configs.Config,
) {
	h := &ViewsHandler{
		usersService: usersService,
		booksService: booksService,
		cfg:          cfg,
	}

	r.Get("/login", h.LoginView)
	r.Get("/register", h.RegisterView)

	middlewares.UseAuthMiddleware(r, cfg, usersService)
	r.Get("/", h.IndexView)
	r.Get("/admin", h.AdminView)
}

func (h *ViewsHandler) IndexView(c *fiber.Ctx) error {
	books, err := h.booksService.GetBooks()
	if err != nil {
		c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.Render("index", fiber.Map{
		"IsAuthenticated": h.usersService.IsAuthenticated(c),
		"Books":           books,
	},
		"layouts/main",
	)
}

func (h *ViewsHandler) AdminView(c *fiber.Ctx) error { // TODO: add admin view
	// books, err := h.booksService.GetBooks()
	// if err != nil {
	// c.Status(fiber.StatusBadRequest).SendString(err.Error())
	// }
	return fmt.Errorf("error")
}

func (h *ViewsHandler) LoginView(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{
		"IsAuthenticated": h.usersService.IsAuthenticated(c),
	},
		"layouts/main",
	)
}

func (h *ViewsHandler) RegisterView(c *fiber.Ctx) error {
	return c.Render("register", fiber.Map{
		"IsAuthenticated": h.usersService.IsAuthenticated(c),
	},
		"layouts/main",
	)
}
