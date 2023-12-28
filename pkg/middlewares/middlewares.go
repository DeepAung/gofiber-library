package middlewares

import (
	"github.com/DeepAung/gofiber-library/modules/users"
	"github.com/gofiber/fiber/v2"
)

type Middleware struct{}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) JwtAuth(usersService *users.UsersService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, err := usersService.VerifyTokenByCookie(c)
		if err != nil {
			usersService.ClearToken(c)
			return c.Redirect("/login")
		}

		return c.Next()
	}
}
