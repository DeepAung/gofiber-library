package viewsHandler

import (
	"strconv"

	"github.com/DeepAung/gofiber-library/configs"
	"github.com/DeepAung/gofiber-library/modules/books/booksService"
	"github.com/DeepAung/gofiber-library/modules/models"
	"github.com/DeepAung/gofiber-library/modules/users/usersService"
	"github.com/DeepAung/gofiber-library/pkg/middlewares"
	"github.com/DeepAung/gofiber-library/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type ViewsHandler struct {
	myerror      *utils.MyError
	usersService *usersService.UsersService
	booksService *booksService.BooksService
	mid          *middlewares.Middleware
	cfg          *configs.Config
}

func NewViewsHandler(
	r fiber.Router,
	myerror *utils.MyError,
	usersService *usersService.UsersService,
	booksService *booksService.BooksService,
	mid *middlewares.Middleware,
	cfg *configs.Config,
) {
	h := &ViewsHandler{
		myerror:      myerror,
		usersService: usersService,
		booksService: booksService,
		mid:          mid,
		cfg:          cfg,
	}

	onlyAuthorized := mid.JwtAccessTokenAuth(usersService)
	onlyUnauthorized := mid.OnlyUnauthorizedAuth(usersService)

	onlyAdmin := mid.OnlyAdmin(usersService)
	setIsAdmin := mid.SetIsAdmin(usersService)
	setOnAdminPage := func(c *fiber.Ctx) error {
		c.Locals("onAdminPage", true)
		return c.Next()
	}

	r.Get("/login", onlyUnauthorized, h.LoginView)
	r.Get("/register", onlyUnauthorized, h.RegisterView)

	r.Get("/", onlyAuthorized, setIsAdmin, h.IndexView)
	r.Get("/books/:id", onlyAuthorized, setIsAdmin, h.DetailView)

	r.Get("/admin", onlyAuthorized, onlyAdmin, setOnAdminPage, h.IndexView)
	r.Get("/admin/books/:id", onlyAuthorized, onlyAdmin, setOnAdminPage, h.DetailView)
	r.Get("/admin/create", onlyAuthorized, onlyAdmin, setOnAdminPage, h.CreateView)
}

func (h *ViewsHandler) IndexView(c *fiber.Ctx) error {
	books, err := h.booksService.GetBooks()

	isAuthenticated := true
	payload := c.Locals("payload").(*models.JwtPayload)
	isAdmin := c.Locals("isAdmin") == true
	onAdminPage := c.Locals("onAdminPage") == true

	if err != nil {
		return h.myerror.RenderErrorPage(
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
	payload := c.Locals("payload").(*models.JwtPayload)
	isAdmin := c.Locals("isAdmin") == true
	onAdminPage := c.Locals("onAdminPage") == true

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return h.myerror.RenderErrorPage(
			c,
			"id should be integer",
			isAuthenticated,
			payload,
			isAdmin,
			onAdminPage,
		)
	}

	book, err := h.booksService.GetBook(id)
	if err != nil {
		return h.myerror.RenderErrorPage(
			c,
			"id should be integer",
			isAuthenticated,
			payload,
			isAdmin,
			onAdminPage,
		)
	}

	isFavorite, err := h.booksService.GetIsFavorite(payload.ID, id)

	return c.Render("detail", fiber.Map{
		"Book":            book,
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
