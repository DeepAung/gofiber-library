package utils

import (
	"github.com/gofiber/fiber/v2"
)

func RenderErrorText(c *fiber.Ctx, msg string) error {
	return c.
		Render("components/error-text", fiber.Map{
			"Error": msg,
		})
}

func RenderErrorPage(
	c *fiber.Ctx,
	msg string,
	isAuthenticated bool,
	payload any,
	isAdmin bool,
	onAdminPage bool,
) error {
	return c.Render("error", fiber.Map{
		"ErrorTitle":      "Error 400",
		"ErrorDetail":     msg,
		"IsAuthenticated": isAuthenticated,
		"Payload":         payload,
		"IsAdmin":         isAdmin,
		"OnAdminPage":     onAdminPage,
	}, "layouts/main")
}
