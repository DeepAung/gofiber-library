package books

import (
	"fmt"

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

func (s *BooksService) GetBooks() (*[]models.Book, error) {
	books := new([]models.Book)
	if err := s.db.Find(books).Error; err != nil {
		return nil, err
	}

	return books, nil
}

func (s *BooksService) GetBook(id int) (*models.Book, error) {
	book := new(models.Book)
	if err := s.db.First(book, id).Error; err != nil {
		return nil, err
	}

	return book, nil
}

func (s *BooksService) CreateBook(book *models.Book) (*models.Book, error) {
	if err := s.db.Create(book).Error; err != nil {
		return nil, err
	}

	return book, nil
}

func (s *BooksService) UpdateBook(book *models.Book, id int) (*models.Book, error) {
	err := s.db.Where("id = ?", id).Updates(book).Error
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (s *BooksService) ToggleFavoriteBook(userId int, bookId int) (bool, error) {
	// user := new(models.User)
	// err := s.db.
	// 	Model(&models.User{}).
	// 	Preload("FavBooks").
	// 	Where("id = ?", userId).
	// 	Find(user).
	// 	Error

	return false, nil
	//
	// var data interface{}
	// s.db.Table("user_favbooks").Where("user_id = ? AND book_id = ?", userId, bookId).Find(data)
	//
	// if err != nil {
	// 	return false, err
	// }
	//
	// has := utils.Has(user.FavBooks, func(book *models.Book) bool {
	// 	return int(book.ID) == bookId
	// })
	//
	// if has {
	// 	err = s.db.Table("user_favbooks")
	// } else {
	//
	// }
	//
	// return false, nil
}

func (s *BooksService) DeleteBook(id int) error {
	result := s.db.Delete(&models.Book{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf(fiber.ErrNotFound.Message)
	}

	return nil
}
