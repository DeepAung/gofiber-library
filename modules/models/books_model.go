package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title         string `json:"title"         gorm:"unique"`
	Author        string `json:"author"`
	Desc          string `json:"description"`
	Content       string `json:"content"`
	FavoriteCount int    `json:"favoriteCount"` // TODO: default zero and fill up other param
}

type CreateBookReq struct {
	Title   string `json:"title"       validate:"required"`
	Author  string `json:"author"      validate:"required"`
	Desc    string `json:"description" validate:"required"`
	Content string `json:"content"     validate:"required"`
}
