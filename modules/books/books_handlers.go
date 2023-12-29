package books

import (
	"fmt"
	"strconv"

	"github.com/DeepAung/gofiber-library/modules/users"
	"github.com/DeepAung/gofiber-library/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type BooksHandler struct {
	validator    *utils.MyValidator
	booksService *BooksService
	usersService *users.UsersService
}

func NewBooksHandler(
	r fiber.Router,
	validator *utils.MyValidator,
	booksService *BooksService,
	usersService *users.UsersService,
) {
	h := &BooksHandler{
		validator:    validator,
		booksService: booksService,
		usersService: usersService,
	}

	r.Post("/books", h.CreateBook)
	r.Post("/books/:id/favorite", h.ToggleFavoriteBook)
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
