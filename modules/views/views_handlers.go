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
	mid          *middlewares.Middleware
	cfg          *configs.Config
}

func NewViewsHandler(
	r fiber.Router,
	usersService *users.UsersService,
	booksService *books.BooksService,
	mid *middlewares.Middleware,
	cfg *configs.Config,
) {
	h := &ViewsHandler{
		usersService: usersService,
		booksService: booksService,
		mid:          mid,
		cfg:          cfg,
	}

	r.Get("/login", h.LoginView)
	r.Get("/register", h.RegisterView)

	r.Use(mid.JwtAuth(usersService))
	r.Get("/", h.IndexView)
	r.Get("/admin", h.AdminView)
}

func (h *ViewsHandler) IndexView(c *fiber.Ctx) error {
	books, err := h.booksService.GetBooks()
	if err != nil {
		c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	payload, err := h.usersService.VerifyTokenByCookie(c)

	return c.Render("index", fiber.Map{
		"IsAuthenticated": err == nil,
		"Payload":         payload,
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
	payload, err := h.usersService.VerifyTokenByCookie(c)

	return c.Render("login", fiber.Map{
		"IsAuthenticated": err == nil,
		"Payload":         payload,
	},
		"layouts/main",
	)
}

func (h *ViewsHandler) RegisterView(c *fiber.Ctx) error {
	payload, err := h.usersService.VerifyTokenByCookie(c)

	return c.Render("register", fiber.Map{
		"IsAuthenticated": err == nil,
		"Payload":         payload,
	},
		"layouts/main",
	)
}
