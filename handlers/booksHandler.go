package handlers

import (
	"strconv"

	"github.com/DeepAung/gofiber-library/pkg/middlewares"
	"github.com/DeepAung/gofiber-library/pkg/utils"
	"github.com/DeepAung/gofiber-library/services"
	"github.com/DeepAung/gofiber-library/types"
	"github.com/gofiber/fiber/v2"
)

type BooksHandler struct {
	booksService *services.BooksService
	usersService *services.UsersService
	mid          *middlewares.Middleware
}

func NewBooksHandler(
	r fiber.Router,
	booksService *services.BooksService,
	usersService *services.UsersService,
	mid *middlewares.Middleware,
) {
	h := &BooksHandler{
		booksService: booksService,
		usersService: usersService,
		mid:          mid,
	}

	onlyAuthorized := mid.JwtAccessTokenAuth(usersService)
	onlyAdmin := mid.OnlyAdmin(usersService)

	r.Post("/books", onlyAuthorized, onlyAdmin, h.CreateBook)
	r.Put("/books/:id", onlyAuthorized, onlyAdmin, h.UpdateBook)
	r.Delete("/books/:id", onlyAuthorized, onlyAdmin, h.DeleteBook)
	r.Post("/books/:id/favorite", onlyAuthorized, h.ToggleFavoriteBook)
}

func (h *BooksHandler) CreateBook(c *fiber.Ctx) error {
	bookReq := new(types.BookReq)
	if err := c.BodyParser(bookReq); err != nil {
		return utils.RenderErrorText(c, err.Error())
	}

	if err := utils.Validate(bookReq); err != nil {
		return utils.RenderErrorText(c, err.Error())
	}

	err := h.booksService.CreateBook(bookReq)
	if err != nil {
		return utils.RenderErrorText(c, err.Error())
	}

	c.Response().Header.Set("HX-Redirect", "/admin")
	return c.SendStatus(fiber.StatusOK)
}

func (h *BooksHandler) UpdateBook(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.RenderErrorText(c, "id should be integer")
	}

	bookReq := new(types.BookReq)
	if err := c.BodyParser(bookReq); err != nil {
		return utils.RenderErrorText(c, err.Error())
	}

	if err := utils.Validate(bookReq); err != nil {
		return utils.RenderErrorText(c, err.Error())
	}

	err = h.booksService.UpdateBook(bookReq, id)
	if err != nil {
		return utils.RenderErrorText(c, err.Error())
	}

	c.Response().Header.Set("HX-Redirect", "/admin")
	return c.SendStatus(fiber.StatusOK)
}

func (h *BooksHandler) DeleteBook(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.RenderErrorText(c, "id should be integer")
	}

	err = h.booksService.DeleteBook(id)
	if err != nil {
		return utils.RenderErrorText(c, err.Error())
	}

	c.Response().Header.Set("HX-Redirect", "/admin")
	return c.SendStatus(fiber.StatusOK)
}

func (h *BooksHandler) ToggleFavoriteBook(c *fiber.Ctx) error {
	bookId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.RenderErrorText(c, "id should be integer")
	}

	book, err := h.booksService.GetBook(bookId)
	if err != nil {
		return utils.RenderErrorText(c, "book not found")
	}

	payload, ok := c.Locals("payload").(*types.JwtPayload)
	if !ok {
		return utils.RenderErrorText(c, "authorization error")
	}

	isFavorite, err := h.booksService.ToggleFavoriteBook(payload.ID, bookId)
	if err != nil {
		return utils.RenderErrorText(c, err.Error())
	}

	if isFavorite {
		book.FavoriteCount += 1
	} else {
		book.FavoriteCount -= 1
	}

	return c.Render("components/star", &fiber.Map{"Book": book, "IsFavorite": isFavorite})
}
