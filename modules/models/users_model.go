package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username"      gorm:"unique"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isAdmin"`
	FavBooks []Book `json:"favoriteBooks" gorm:"many2many:user_favbooks"`
}

type LoginReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterReq struct {
	Username  string `json:"username"  validate:"required"`
	Password  string `json:"password"  validate:"required"`
	Password2 string `json:"password2" validate:"required"`
}
