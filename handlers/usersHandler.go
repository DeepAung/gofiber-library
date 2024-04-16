package handlers

import (
	"github.com/DeepAung/gofiber-library/pkg/middlewares"
	"github.com/DeepAung/gofiber-library/pkg/utils"
	"github.com/DeepAung/gofiber-library/services"
	"github.com/DeepAung/gofiber-library/types"
	"github.com/gofiber/fiber/v2"
)

type UsersHandler struct {
	usersSvc *services.UsersService
	mid      *middlewares.Middleware
}

func NewUsersHandler(usersSvc *services.UsersService, mid *middlewares.Middleware) *UsersHandler {
	return &UsersHandler{
		usersSvc: usersSvc,
		mid:      mid,
	}
}

func (h *UsersHandler) Login(c *fiber.Ctx) error {
	loginReq := new(types.LoginReq)
	if err := c.BodyParser(loginReq); err != nil {
		return utils.RenderErrorText(c, err.Error())
	}

	if err := utils.Validate(loginReq); err != nil {
		return utils.RenderErrorText(c, err.Error())
	}

	err := h.usersSvc.Login(loginReq, c)
	if err != nil {
		return utils.RenderErrorText(c, err.Error())
	}

	c.Response().Header.Set("HX-Redirect", "/")
	return c.SendStatus(fiber.StatusOK)
}

func (h *UsersHandler) Register(c *fiber.Ctx) error {
	registerReq := new(types.RegisterReq)
	if err := c.BodyParser(registerReq); err != nil {
		return utils.RenderErrorText(c, err.Error())
	}

	if err := utils.Validate(registerReq); err != nil {
		return utils.RenderErrorText(c, err.Error())
	}

	err := h.usersSvc.Register(registerReq)
	if err != nil {
		return utils.RenderErrorText(c, err.Error())
	}

	c.Response().Header.Set("HX-Redirect", "/login")
	return c.SendStatus(fiber.StatusOK)
}

func (h *UsersHandler) Logout(c *fiber.Ctx) error {
	h.usersSvc.ClearToken(c)

	c.Response().Header.Set("HX-Redirect", "/login")
	return c.SendStatus(fiber.StatusOK)
}

func (h *UsersHandler) UpdateTokens(c *fiber.Ctx) error {
	_, err := h.usersSvc.UpdateTokens(c)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
