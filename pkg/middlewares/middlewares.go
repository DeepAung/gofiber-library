package middlewares

import (
	"github.com/DeepAung/gofiber-library/modules/users/usersService"
	"github.com/gofiber/fiber/v2"
)

type Middleware struct{}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) JwtAccessTokenAuth(usersService *usersService.UsersService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		payload, err := usersService.VerifyTokenByCookie(c, "access_token")

		if err != nil {
			err = usersService.UpdateTokens(c)
		}

		if err != nil {
			usersService.ClearToken(c)
			return c.Redirect("/login")
		} else {
		}

		c.Locals("payload", payload)
		return c.Next()
	}
}

func (m *Middleware) JwtRefreshTokenAuth(usersService *usersService.UsersService) fiber.Handler {
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

func (m *Middleware) OnlyUnauthorizedAuth(usersService *usersService.UsersService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, err := usersService.VerifyTokenByCookie(c, "access_token")

		if err != nil {
			return c.Next()
		}

		return c.Redirect("/")
	}
}
