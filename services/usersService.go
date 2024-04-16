package services

import (
	"fmt"
	"time"

	"github.com/DeepAung/gofiber-library/pkg/configs"
	"github.com/DeepAung/gofiber-library/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UsersService struct {
	db  *gorm.DB
	cfg *configs.Config
}

func NewUsersService(db *gorm.DB, cfg *configs.Config) *UsersService {
	return &UsersService{
		db:  db,
		cfg: cfg,
	}
}

const (
	AccessTokenExpTime  = 15 * time.Minute
	RefreshTokenExpTime = 7 * 24 * time.Hour
)

func (s *UsersService) Login(req *types.LoginReq, c *fiber.Ctx) error {
	user := new(types.User)
	err := s.db.
		Model(&types.User{}).
		Where("username = ?", req.Username).
		First(user).Error
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if !s.checkPassword(req.Password, user.Password) {
		return fmt.Errorf("password incorrect")
	}

	accessToken, err := s.generateToken(int(user.ID), user.Username, AccessTokenExpTime)
	if err != nil {
		return err
	}

	refreshToken, err := s.generateToken(int(user.ID), user.Username, RefreshTokenExpTime)
	if err != nil {
		return err
	}

	err = s.db.
		Model(&types.User{}).
		Where("id = ?", user.ID).
		Update("refresh_token", refreshToken).Error
	if err != nil {
		return err
	}

	s.setCookie(c, "access_token", accessToken, AccessTokenExpTime)
	s.setCookie(c, "refresh_token", refreshToken, RefreshTokenExpTime)

	return nil
}

func (s *UsersService) Register(req *types.RegisterReq) error {
	if req.Password != req.Password2 {
		return fmt.Errorf("password is not the same")
	}

	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return err
	}

	user := types.User{
		Username: req.Username,
		Password: hashedPassword,
	}
	err = s.db.Create(&user).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *UsersService) UpdateTokens(c *fiber.Ctx) (*types.JwtPayload, error) {
	cookieRefreshToken := c.Cookies("refresh_token")
	if cookieRefreshToken == "" {
		return nil, fmt.Errorf("refresh token not found")
	}

	payload, err := s.VerifyToken(cookieRefreshToken)
	if err != nil {
		return nil, err
	}

	dbRefreshToken := ""
	err = s.db.Model(&types.User{}).
		Where("id = ?", payload.ID).
		Select("refresh_token").
		First(&dbRefreshToken).
		Error
	if err != nil {
		return nil, err
	}

	if cookieRefreshToken != dbRefreshToken {
		return nil, fmt.Errorf("incorrect refresh token")
	}

	newAccessToken, err := s.generateToken(
		payload.ID,
		payload.Username,
		AccessTokenExpTime,
	)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.generateToken(
		payload.ID,
		payload.Username,
		RefreshTokenExpTime,
	)
	if err != nil {
		return nil, err
	}

	err = s.db.
		Model(&types.User{}).
		Where("id = ?", payload.ID).
		Update("refresh_token", newRefreshToken).
		Error
	if err != nil {
		return nil, err
	}

	s.setCookie(c, "access_token", newAccessToken, AccessTokenExpTime)
	s.setCookie(c, "refresh_token", newRefreshToken, RefreshTokenExpTime)

	return payload, nil
}

func (s *UsersService) VerifyToken(tokenString string) (*types.JwtPayload, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("token not found")
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&types.JwtClaim{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(s.cfg.JwtSecret), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*types.JwtClaim)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return &claims.Payload, nil
}

func (s *UsersService) ClearToken(c *fiber.Ctx) {
	s.setCookie(c, "access_token", "deleted", -2*time.Hour)
	s.setCookie(c, "refresh_token", "deleted", -2*time.Hour)
}

func (s *UsersService) IsAdmin(id int) (bool, error) {
	var isAdmin bool
	err := s.db.Model(&types.User{}).Where("id = ?", id).Select("is_admin").First(&isAdmin).Error
	if err != nil {
		return false, err
	}

	return isAdmin, nil
}

func (s *UsersService) setCookie(
	c *fiber.Ctx,
	name string,
	value string,
	expTime time.Duration,
) {
	c.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    value,
		Expires:  time.Now().Add(expTime),
		HTTPOnly: true,
		Secure:   true,
		Path:     "/",
	})
}

func (s *UsersService) generateToken(
	userId int,
	username string,
	expTime time.Duration,
) (string, error) {
	claims := types.JwtClaim{
		Payload: types.JwtPayload{
			ID:       userId,
			Username: username,
			Exp:      time.Now().Add(expTime).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	return token.SignedString([]byte(s.cfg.JwtSecret))
}

func (s *UsersService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *UsersService) checkPassword(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
