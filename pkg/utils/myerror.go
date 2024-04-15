package utils

import (
	"github.com/DeepAung/gofiber-library/types"
	"github.com/gofiber/fiber/v2"
)

type MyError struct{}

func NewMyError() *MyError {
	return &MyError{}
}

func (e *MyError) SendErrorText(c *fiber.Ctx, msg string) error {
	return c.
		Status(fiber.StatusBadRequest).
		Render("components/error-text", fiber.Map{
			"Error": msg,
		})
}

func (e *MyError) RenderErrorPage(
	c *fiber.Ctx,
	msg string,
	isAuthenticated bool,
	payload *types.JwtPayload,
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
