package dummy

import (
	"fmt"

	"github.com/DeepAung/gofiber-library/modules/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type DummyHandler struct {
	db *gorm.DB
}

func NewDummyHandler(r fiber.Router, db *gorm.DB) {
	r.Post("/books", func(c *fiber.Ctx) error {
		book := new(models.Book)
		if err := c.BodyParser(book); err != nil {
			return err
		}

		fmt.Println("dummy book: ", book)

		if err := db.Create(book).Error; err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusCreated)
	})
}
