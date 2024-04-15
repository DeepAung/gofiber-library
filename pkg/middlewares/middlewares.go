package middlewares

import (
	"github.com/DeepAung/gofiber-library/services"
	"github.com/DeepAung/gofiber-library/types"
	"github.com/gofiber/fiber/v2"
)

type Middleware struct{}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) PageNotFound(usersService *services.UsersService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		payload, err := usersService.VerifyTokenByCookie(c, "access_token")

		return c.Render("error", fiber.Map{
			"IsAuthenticated": err == nil,
			"Payload":         payload,
			"ErrorTitle":      "Error 404",
			"ErrorDetail":     "This page is not found",
		},
			"layouts/main",
		)
	}
}

func (m *Middleware) SetIsAdmin(usersService *services.UsersService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		payload, ok := c.Locals("payload").(*types.JwtPayload)
		if !ok {
			usersService.ClearToken(c)
			return c.Redirect("/login")
		}

		isAdmin, err := usersService.IsAdmin(payload.ID)
		if err != nil {
			return c.Redirect("/")
		}

		c.Locals("isAdmin", isAdmin)

		return c.Next()
	}
}

func (m *Middleware) OnlyAdmin(usersService *services.UsersService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		payload, ok := c.Locals("payload").(*types.JwtPayload)
		if !ok {
			usersService.ClearToken(c)
			return c.Redirect("/login")
		}

		isAdmin, err := usersService.IsAdmin(payload.ID)
		if err != nil {
			return c.Redirect("/")
		}

		c.Locals("isAdmin", isAdmin)
		if isAdmin {
			return c.Next()
		} else {
			return c.Redirect("/")
		}
	}
}

func (m *Middleware) JwtAccessTokenAuth(usersService *services.UsersService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		payload, err := usersService.VerifyTokenByCookie(c, "access_token")
		if err == nil {
			c.Locals("payload", payload)
			return c.Next()
		}

		payload, err = usersService.UpdateTokens(c)
		if err != nil {
			usersService.ClearToken(c)
			return c.Redirect("/login")
		}

		c.Locals("payload", payload)
		return c.Next()
	}
}

func (m *Middleware) JwtRefreshTokenAuth(usersService *services.UsersService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		payload, err := usersService.VerifyTokenByCookie(c, "refresh_token")
		if err != nil {
			usersService.ClearToken(c)
			return c.Redirect("/login")
		}

		c.Locals("payload", payload)
		return c.Next()
	}
}

func (m *Middleware) OnlyUnauthorizedAuth(usersService *services.UsersService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, err := usersService.VerifyTokenByCookie(c, "access_token")
		if err != nil {
			return c.Next()
		}

		return c.Redirect("/")
	}
}
