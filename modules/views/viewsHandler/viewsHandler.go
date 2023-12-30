package viewsHandler

import (
	"fmt"
	"strconv"

	"github.com/DeepAung/gofiber-library/configs"
	"github.com/DeepAung/gofiber-library/modules/books/booksService"
	"github.com/DeepAung/gofiber-library/modules/models"
	"github.com/DeepAung/gofiber-library/modules/users/usersService"
	"github.com/DeepAung/gofiber-library/pkg/middlewares"
	"github.com/gofiber/fiber/v2"
)

type ViewsHandler struct {
	usersService *usersService.UsersService
	booksService *booksService.BooksService
	mid          *middlewares.Middleware
	cfg          *configs.Config
}

func NewViewsHandler(
	r fiber.Router,
	usersService *usersService.UsersService,
	booksService *booksService.BooksService,
	mid *middlewares.Middleware,
	cfg *configs.Config,
) {
	h := &ViewsHandler{
		usersService: usersService,
		booksService: booksService,
		mid:          mid,
		cfg:          cfg,
	}

	onlyAuthorized := mid.JwtAccessTokenAuth(usersService)
	onlyUnauthorized := mid.OnlyUnauthorizedAuth(usersService)

	r.Get("/login", onlyUnauthorized, h.LoginView)
	r.Get("/register", onlyUnauthorized, h.RegisterView)

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

	payload, ok := c.Locals("payload").(*models.JwtPayload)
	if !ok {
		return fmt.Errorf("authorization error")
	}

	book, err := h.booksService.GetBook(id)
	if err != nil {
		return err
	}

	isFavorite, err := h.booksService.GetIsFavorite(payload.ID, id)

	return c.Render("detail", fiber.Map{
		"IsAuthenticated": true,
		"Payload":         payload,
		"Book":            book,
		"IsFavorite":      isFavorite,
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
