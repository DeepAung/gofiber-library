package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title         string `json:"title"         gorm:"unique" validate:"required"`
	Author        string `json:"author"                      validate:"required"`
	Desc          string `json:"desc"                        validate:"required"`
	FavoriteCount int    `json:"favoriteCount"` // TODO: default zero and fill up other param
}

type BookReq struct {
}
