package types

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title         string `json:"title"         form:"title"         gorm:"unique"`
	Author        string `json:"author"        form:"author"`
	Desc          string `json:"description"   form:"description"`
	Content       string `json:"content"       form:"content"`
	FavoriteCount int    `json:"favoriteCount" form:"favoriteCount" gorm:"default:0"`
}

type BookReq struct {
	Title   string `json:"title"       form:"title"       validate:"required"`
	Author  string `json:"author"      form:"author"      validate:"required"`
	Desc    string `json:"description" form:"description"`
	Content string `json:"content"     form:"content"`
}
