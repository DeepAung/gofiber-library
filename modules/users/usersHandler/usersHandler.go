package usersHandler

import (
	"github.com/DeepAung/gofiber-library/modules/models"
	"github.com/DeepAung/gofiber-library/modules/users/usersService"
	"github.com/DeepAung/gofiber-library/pkg/middlewares"
	"github.com/DeepAung/gofiber-library/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type UsersHandler struct {
	validator    *utils.MyValidator
	usersService *usersService.UsersService
	mid          *middlewares.Middleware
}

func NewUsersHandler(
	r fiber.Router,
	validator *utils.MyValidator,
	usersService *usersService.UsersService,
	mid *middlewares.Middleware,
) {
	h := &UsersHandler{
		validator:    validator,
		usersService: usersService,
		mid:          mid,
	}

	onlyAuthorized := mid.JwtAccessTokenAuth(usersService)
	onlyUnauthorized := mid.OnlyUnauthorizedAuth(usersService)

	r.Post("/login", onlyUnauthorized, h.Login)
	r.Post("/register", onlyUnauthorized, h.Register)

	r.Post("/logout", onlyAuthorized, h.Logout)
	r.Post("/refresh", onlyAuthorized, h.UpdateTokens)
}

func (h *UsersHandler) Login(c *fiber.Ctx) error {
	loginReq := new(models.LoginReq)
	if err := c.BodyParser(loginReq); err != nil {
		return err // TODO: return error components
	}

	if err := h.validator.Validate(loginReq); err != nil {
		return err
	}

	err := h.usersService.Login(loginReq, c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			Render("components/auth-error", fiber.Map{"Errors": []string{err.Error()}})
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

	err := h.usersService.Register(registerReq)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			Render("components/auth-error", fiber.Map{"Errors": []string{err.Error()}})
	}

	c.Response().Header.Set("HX-Redirect", "/login")
	return c.SendStatus(fiber.StatusOK)
}

func (h *UsersHandler) Logout(c *fiber.Ctx) error {
	h.usersService.ClearToken(c)

	c.Response().Header.Set("HX-Redirect", "/login")
	return c.SendStatus(fiber.StatusOK)
}

func (h *UsersHandler) UpdateTokens(c *fiber.Ctx) error {
	err := h.usersService.UpdateTokens(c)
	if err != nil {
		log.Error("ERROR: ", err)
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
