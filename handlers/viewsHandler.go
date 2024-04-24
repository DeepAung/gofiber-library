package handlers

import (
	"html/template"
	"strconv"

	"github.com/DeepAung/gofiber-library/pkg/configs"
	"github.com/DeepAung/gofiber-library/pkg/middlewares"
	"github.com/DeepAung/gofiber-library/pkg/utils"
	"github.com/DeepAung/gofiber-library/services"
	"github.com/DeepAung/gofiber-library/types"
	"github.com/gofiber/fiber/v2"
)

type ViewsHandler struct {
	booksSvc *services.BooksService
	mid      *middlewares.Middleware
	cfg      *configs.Config
}

func NewViewsHandler(
	booksSvc *services.BooksService,
	mid *middlewares.Middleware,
	cfg *configs.Config,
) *ViewsHandler {
	return &ViewsHandler{
		booksSvc: booksSvc,
		mid:      mid,
		cfg:      cfg,
	}
}

func (h *ViewsHandler) IndexView(c *fiber.Ctx) error {
	books, err := h.booksSvc.GetBooks()

	isAuthenticated := true
	payload := c.Locals("payload").(*types.JwtPayload)
	isAdmin := c.Locals("isAdmin") == true
	onAdminPage := c.Locals("onAdminPage") == true

	if err != nil {
		return utils.RenderErrorPage(
			c,
			err.Error(),
			isAuthenticated,
			payload,
			isAdmin,
			onAdminPage,
		)
	}

	return c.Render("index", fiber.Map{
		"Books":           books,
		"IsAuthenticated": isAuthenticated,
		"Payload":         payload,
		"IsAdmin":         isAdmin,
		"OnAdminPage":     onAdminPage,
	}, "layouts/main")
}

func (h *ViewsHandler) DetailView(c *fiber.Ctx) error {
	isAuthenticated := true
	payload := c.Locals("payload").(*types.JwtPayload)
	isAdmin := c.Locals("isAdmin") == true
	onAdminPage := c.Locals("onAdminPage") == true

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.RenderErrorPage(
			c,
			"id should be integer",
			isAuthenticated,
			payload,
			isAdmin,
			onAdminPage,
		)
	}

	book, err := h.booksSvc.GetBook(id)
	if err != nil {
		return utils.RenderErrorPage(
			c,
			"id should be integer",
			isAuthenticated,
			payload,
			isAdmin,
			onAdminPage,
		)
	}

	isFavorite, err := h.booksSvc.GetIsFavorite(payload.ID, id)
	if err != nil {
		return utils.RenderErrorPage(
			c,
			"get is favorite failed",
			isAuthenticated,
			payload,
			isAdmin,
			onAdminPage,
		)
	}

	return c.Render("detail", fiber.Map{
		"Book":            book,
		"Content":         template.HTML(book.Content),
		"IsFavorite":      isFavorite,
		"IsAuthenticated": true,
		"Payload":         payload,
		"IsAdmin":         isAdmin,
		"OnAdminPage":     onAdminPage,
	},
		"layouts/main",
	)
}

func (h *ViewsHandler) CreateView(c *fiber.Ctx) error {
	return c.Render("create", fiber.Map{
		"IsAuthenticated": true,
		"Payload":         c.Locals("payload"),
		"IsAdmin":         true,
		"OnAdminPage":     true,
	}, "layouts/main")
}

func (h *ViewsHandler) LoginView(c *fiber.Ctx) error {
	return c.Render("login", nil, "layouts/main")
}

func (h *ViewsHandler) RegisterView(c *fiber.Ctx) error {
	return c.Render("register", nil, "layouts/main")
}
