package middlewares

import (
	"github.com/DeepAung/gofiber-library/configs"
	"github.com/DeepAung/gofiber-library/modules/users"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func UseAuthMiddleware(r fiber.Router, cfg *configs.Config, s *users.UsersService) {
	r.Use(CookieToBearer)
	r.Use(jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(cfg.JwtSecret)},
		ErrorHandler: JwtErrorHandler(s),
	}))
}

func CookieToBearer(c *fiber.Ctx) error {
	accessToken := c.Cookies("access_token")
	c.Request().Header.Set("Authorization", "Bearer "+accessToken)

	return c.Next()
}

func JwtErrorHandler(s *users.UsersService) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		s.ClearToken(c) // TODO: shouldn't do this
		return c.Redirect("/login")
	}
}
