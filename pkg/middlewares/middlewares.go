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

func (m *Middleware) PageNotFound(usersSvc *services.UsersService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		payload, err := usersSvc.VerifyToken(c.Cookies("access_token"))

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

func (m *Middleware) JwtAccessTokenAuth(usersSvc *services.UsersService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		payload, err := usersSvc.VerifyToken(c.Cookies("access_token"))
		if err != nil {
			if payload, err = usersSvc.UpdateTokens(c); err != nil {
				usersSvc.ClearToken(c)
				return c.Redirect("/login")
			}
		}

		c.Locals("payload", payload)
		return c.Next()
	}
}

func (m *Middleware) JwtRefreshTokenAuth(usersSvc *services.UsersService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		payload, err := usersSvc.VerifyToken(c.Cookies("refresh_token"))
		if err != nil {
			usersSvc.ClearToken(c)
			return c.Redirect("/login")
		}

		c.Locals("payload", payload)
		return c.Next()
	}
}

func (m *Middleware) OnlyUnauthorizedAuth(usersSvc *services.UsersService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, err := usersSvc.VerifyToken(c.Cookies("access_token"))
		if err != nil {
			return c.Next()
		}

		return c.Redirect("/")
	}
}

func (m *Middleware) OnlyAdmin(usersSvc *services.UsersService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		payload, ok := c.Locals("payload").(*types.JwtPayload)
		if !ok {
			usersSvc.ClearToken(c)
			return c.Redirect("/login")
		}

		isAdmin, err := usersSvc.IsAdmin(payload.ID)
		if err != nil {
			usersSvc.ClearToken(c)
			return c.Redirect("/login")
		}

		c.Locals("isAdmin", isAdmin)
		if isAdmin {
			return c.Next()
		} else {
			return c.Redirect("/")
		}
	}
}

func (m *Middleware) SetIsAdmin(usersSvc *services.UsersService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		payload, ok := c.Locals("payload").(*types.JwtPayload)
		if !ok {
			usersSvc.ClearToken(c)
			return c.Redirect("/login")
		}

		isAdmin, err := usersSvc.IsAdmin(payload.ID)
		if err != nil {
			usersSvc.ClearToken(c)
			return c.Redirect("/login")
		}

		c.Locals("isAdmin", isAdmin)

		return c.Next()
	}
}
