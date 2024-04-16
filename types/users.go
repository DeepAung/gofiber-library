package types

import (
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `json:"username"      form:"username"      gorm:"unique"`
	Password     string `json:"password"      form:"password"`
	RefreshToken string `json:"-"             form:"-"`
	IsAdmin      bool   `json:"isAdmin"       form:"isAdmin"`
	FavBooks     []Book `json:"favoriteBooks" form:"favoriteBooks" gorm:"many2many:user_favbooks"`
}

type UserFavbooks struct {
	UserID int `gorm:"primaryKey"`
	BookID int `gorm:"primaryKey"`
}

type LoginReq struct {
	Username string `json:"username" form:"username" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

type RegisterReq struct {
	Username  string `json:"username"  form:"username"  validate:"required"`
	Password  string `json:"password"  form:"password"  validate:"required"`
	Password2 string `json:"password2" form:"password2" validate:"required"`
}

type JwtClaim struct {
	Payload JwtPayload
	jwt.RegisteredClaims
}

type JwtPayload struct {
	ID       int    `json:"id"       form:"id"`
	Username string `json:"username" form:"username"`
	Exp      int64  `json:"exp"      form:"exp"`
}
