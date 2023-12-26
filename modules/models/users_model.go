package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string         `json:"username" gorm:"unique"`
	Password string         `json:"password"`
	Books    []BorrowRecord `json:"books"`
}

type LoginReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginRes struct {
	Token string `json:"token"`
}

type RegisterReq struct {
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"`
	Password2 string `json:"password2" validate:"required"`
}
