package users

import (
	"github.com/DeepAung/gofiber-library/modules/models"
	"github.com/DeepAung/gofiber-library/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type UsersHandler struct {
	validator *utils.MyValidator
	service   *UsersService
}

func NewUsersHandler(r fiber.Router, validator *utils.MyValidator, service *UsersService) {
	h := &UsersHandler{
		validator: validator,
		service:   service,
	}

	r.Post("/login", h.Login)
	r.Post("/register", h.Register)
	r.Post("/logout", h.Logout)
}

func (h *UsersHandler) Login(c *fiber.Ctx) error {
	loginReq := new(models.LoginReq)
	if err := c.BodyParser(loginReq); err != nil {
		return err // TODO: return error components
	}

	if err := h.validator.Validate(loginReq); err != nil {
		return err
	}

	err := h.service.Login(loginReq, c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			Render("components/error", fiber.Map{"Errors": []string{err.Error()}})
	}

	c.Response().Header.Set("HX-Redirect", "/")
	return c.SendStatus(fiber.StatusOK)
}

func (h *UsersHandler) Register(c *fiber.Ctx) error {
	registerReq := new(models.RegisterReq)
	if err := c.BodyParser(registerReq); err != nil {
		return err
	}

	if err := h.validator.Validate(registerReq); err != nil {
		return err
	}

	err := h.service.Register(registerReq)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			Render("components/error", fiber.Map{"Errors": []string{err.Error()}})
	}

	c.Response().Header.Set("HX-Redirect", "/login")
	return c.SendStatus(fiber.StatusOK)
}

func (h *UsersHandler) Logout(c *fiber.Ctx) error { // TODO: working on clearing cookies
	h.service.ClearToken(c)
	println("cookies: ", c.Cookies("access_token"), " | ", c.Cookies("refresh_token"))

	c.Response().Header.Set("HX-Redirect", "/login")
	return c.SendStatus(fiber.StatusOK)
}
