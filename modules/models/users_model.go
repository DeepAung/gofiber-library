package models

import (
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `json:"username"      gorm:"unique"`
	Password     string `json:"password"`
	RefreshToken string `json:"-"`
	IsAdmin      bool   `json:"isAdmin"`
	FavBooks     []Book `json:"favoriteBooks" gorm:"many2many:user_favbooks"`
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

type JwtClaim struct {
	Payload JwtPayload
	jwt.RegisteredClaims
}

type JwtPayload struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Exp      int64  `json:"exp"`
}
