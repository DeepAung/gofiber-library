package books

import (
	"fmt"
	"strconv"

	"github.com/DeepAung/gofiber-library/db"
	"github.com/DeepAung/gofiber-library/modules/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type BooksService struct {
	db *gorm.DB
}

func NewBooksService(db *gorm.DB) *BooksService {
	return &BooksService{
		db: db,
	}
}

func getBooks() (*[]models.Book, error) {
	books := new([]models.Book)
	if err := db.DB.Find(books).Error; err != nil {
		return nil, err
	}

	return books, nil
}

func getBook(id int) (*models.Book, error) {
	// id, err := strconv.Atoi(c.Params("id"))
	// if err != nil {
	// 	return nil, fmt.Errorf("ID should be integer")
	// }

	book := new(models.Book)
	if err := db.DB.First(book, id).Error; err != nil {
		return nil, err
	}

	return book, nil
}

func createBook(bookReq *models.BookReq) (*models.Book, error) {
	var book models.Book
	// c.BodyParser(&book)

	// if err := validator.Validator.Struct(book); err != nil {
	// 	return nil, err
	// }

	if err := c.BodyParser(&book); err != nil {
		return nil, err
	}

	err := db.DB.Create(&book).Error
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func updateBook(c *fiber.Ctx) (*models.Book, error) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return nil, fmt.Errorf("ID should be integer")
	}

	var book models.Book

	if err := c.BodyParser(&book); err != nil {
		return nil, err
	}

	err = db.DB.Where("id = ?", id).Updates(&book).Error
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func deleteBook(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fmt.Errorf("ID should be integer")
	}

	result := db.DB.Delete(&models.Book{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf(fiber.ErrNotFound.Message)
	}

	return nil
}
