package booksHandler

import (
	"fmt"
	"strconv"

	"github.com/DeepAung/gofiber-library/modules/books/booksService"
	"github.com/DeepAung/gofiber-library/modules/users/usersService"
	"github.com/DeepAung/gofiber-library/pkg/middlewares"
	"github.com/DeepAung/gofiber-library/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type BooksHandler struct {
	validator    *utils.MyValidator
	booksService *booksService.BooksService
	usersService *usersService.UsersService
	mid          *middlewares.Middleware
}

func NewBooksHandler(
	r fiber.Router,
	validator *utils.MyValidator,
	booksService *booksService.BooksService,
	usersService *usersService.UsersService,
	mid *middlewares.Middleware,
) {
	h := &BooksHandler{
		validator:    validator,
		booksService: booksService,
		usersService: usersService,
		mid:          mid,
	}

	onlyAuthorized := mid.JwtAccessTokenAuth(usersService)

	r.Post("/books", onlyAuthorized, h.CreateBook)
	r.Post("/books/:id/favorite", onlyAuthorized, h.ToggleFavoriteBook)
}

func (h *BooksHandler) CreateBook(c *fiber.Ctx) error {
	return fmt.Errorf("error")
}

func (h *BooksHandler) ToggleFavoriteBook(c *fiber.Ctx) error {
	println("testtttttttttttttttt")
	bookId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fmt.Errorf("id should be integer")
	}

	payload, err := h.usersService.VerifyTokenByCookie(c, "access_token")
	if err != nil {
		return err
	}
	userId := payload.ID

	toggled, err := h.booksService.ToggleFavoriteBook(userId, bookId)
	if err != nil {
		return err
	}

	return c.Render("components/star", &fiber.Map{
		"Toggled": toggled,
	},
		"layouts/main",
	)
}
