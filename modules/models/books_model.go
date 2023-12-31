package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title         string `json:"title"         gorm:"unique"`
	Author        string `json:"author"`
	Desc          string `json:"description"`
	Content       string `json:"content"`
	FavoriteCount int    `json:"favoriteCount" gorm:"default:0"`
}

type BookReq struct {
	Title   string `json:"title"       validate:"required"`
	Author  string `json:"author"      validate:"required"`
	Desc    string `json:"description"`
	Content string `json:"content"`
}
