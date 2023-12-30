package booksHandler

import (
	"fmt"
	"strconv"

	"github.com/DeepAung/gofiber-library/modules/books/booksService"
	"github.com/DeepAung/gofiber-library/modules/models"
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
	bookId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fmt.Errorf("id should be integer")
	}

	book, err := h.booksService.GetBook(bookId)
	if err != nil {
		return fmt.Errorf("book not found")
	}

	payload, ok := c.Locals("payload").(*models.JwtPayload)
	if !ok {
		return fmt.Errorf("authorization error")
	}

	isFavorite, err := h.booksService.ToggleFavoriteBook(payload.ID, bookId)
	if err != nil {
		return err
	}

	if isFavorite {
		book.FavoriteCount += 1
	} else {
		book.FavoriteCount -= 1
	}

	return c.Render("components/star", &fiber.Map{"Book": book, "IsFavorite": isFavorite})
}
