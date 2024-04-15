package services

import (
	"errors"
	"fmt"

	"github.com/DeepAung/gofiber-library/types"
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

func (s *BooksService) GetBooks() (*[]types.Book, error) {
	books := new([]types.Book)
	if err := s.db.Find(books).Error; err != nil {
		return nil, err
	}

	return books, nil
}

func (s *BooksService) GetBook(id int) (*types.Book, error) {
	book := new(types.Book)
	if err := s.db.First(book, id).Error; err != nil {
		return nil, err
	}

	return book, nil
}

func (s *BooksService) CreateBook(bookReq *types.BookReq) error {
	book := &types.Book{
		Title:         bookReq.Title,
		Author:        bookReq.Author,
		Desc:          bookReq.Desc,
		Content:       bookReq.Content,
		FavoriteCount: 0,
	}

	return s.db.Create(book).Error
}

func (s *BooksService) UpdateBook(bookReq *types.BookReq, id int) error {
	book := &types.Book{
		Title:   bookReq.Title,
		Author:  bookReq.Author,
		Desc:    bookReq.Desc,
		Content: bookReq.Content,
	}

	return s.db.Where("id = ?", id).Updates(book).Error
}

func (s *BooksService) DeleteBook(id int) error {
	result := s.db.Delete(&types.Book{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf(fiber.ErrNotFound.Message)
	}

	return nil
}

func (s *BooksService) GetIsFavorite(userId int, bookId int) (bool, error) {
	err := s.db.First(&types.UserFavbooks{UserID: userId, BookID: bookId}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (s *BooksService) ToggleFavoriteBook(userId int, bookId int) (bool, error) {
	isFavorite, err := s.GetIsFavorite(userId, bookId)
	if err != nil {
		return false, err
	}

	if isFavorite {
		return s.UnfavoriteBook(userId, bookId)
	} else {
		return s.FavoriteBook(userId, bookId)
	}
}

func (s *BooksService) FavoriteBook(userId int, bookId int) (bool, error) {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&types.UserFavbooks{UserID: userId, BookID: bookId}).Error
		if err != nil {
			return err
		}

		err = tx.Model(&types.Book{}).
			Where("id = ?", bookId).
			Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error
		if err != nil {
			return err
		}

		return nil
	})

	return true, err
}

func (s *BooksService) UnfavoriteBook(userId int, bookId int) (bool, error) {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Delete(&types.UserFavbooks{UserID: userId, BookID: bookId}).Error
		if err != nil {
			return err
		}

		err = tx.Model(&types.Book{}).
			Where("id = ?", bookId).
			Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).
			Error
		if err != nil {
			return err
		}

		return nil
	})

	return false, err
}
