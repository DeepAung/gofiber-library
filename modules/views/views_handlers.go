package views

import (
	"fmt"
	"strconv"

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

	onlyUnauthorized := mid.OnlyUnauthorizedAuth(usersService)
	r.Get("/login", onlyUnauthorized, h.LoginView)
	r.Get("/register", onlyUnauthorized, h.RegisterView)

	onlyAuthorized := mid.JwtAccessTokenAuth(usersService)
	r.Get("/", onlyAuthorized, h.IndexView)
	r.Get("/books/:id", onlyAuthorized, h.DetailView)
	r.Get("/admin", onlyAuthorized, h.AdminView)
}

func (h *ViewsHandler) IndexView(c *fiber.Ctx) error {
	books, err := h.booksService.GetBooks()
	if err != nil {
		c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.Render("index", fiber.Map{
		"IsAuthenticated": true,
		"Payload":         c.Locals("payload"),
		"Books":           books,
	},
		"layouts/main",
	)
}

func (h *ViewsHandler) DetailView(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fmt.Errorf("id should be integer")
	}

	book, err := h.booksService.GetBook(id)
	if err != nil {
		return err
	}

	return c.Render("detail", fiber.Map{
		"IsAuthenticated": true,
		"Payload":         c.Locals("payload"),
		"Book":            book,
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
		"IsAuthenticated": false,
		"Payload":         nil,
	},
		"layouts/main",
	)
}

func (h *ViewsHandler) RegisterView(c *fiber.Ctx) error {
	return c.Render("register", fiber.Map{
		"IsAuthenticated": false,
		"Payload":         nil,
	},
		"layouts/main",
	)
}
