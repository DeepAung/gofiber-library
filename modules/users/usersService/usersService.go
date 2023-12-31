package usersService

import (
	"fmt"
	"time"

	"github.com/DeepAung/gofiber-library/configs"
	"github.com/DeepAung/gofiber-library/modules/models"
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

const AccessTokenExpTime = 15 * time.Minute
const RefreshTokenExpTime = 7 * 24 * time.Hour

func (s *UsersService) Login(loginReq *models.LoginReq, c *fiber.Ctx) error {
	user := new(models.User)
	err := s.db.
		Where(&models.User{Username: loginReq.Username}).
		First(user).Error
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if !s.CheckPassword(loginReq.Password, user.Password) {
		return fmt.Errorf("password incorrect")
	}

	accessToken, err := s.GenerateToken(
		int(user.ID),
		user.Username,
		AccessTokenExpTime,
	)
	if err != nil {
		return err
	}

	refreshToken, err := s.GenerateToken(
		int(user.ID),
		user.Username,
		RefreshTokenExpTime,
	)
	if err != nil {
		return err
	}

	err = s.db.
		Model(&models.User{}).
		Where("id = ?", user.ID).
		Update("refresh_token", refreshToken).
		Error
	if err != nil {
		return err
	}

	s.SetCookie(c, "access_token", accessToken, AccessTokenExpTime)
	s.SetCookie(c, "refresh_token", refreshToken, RefreshTokenExpTime)

	return nil
}

func (s *UsersService) Register(registerReq *models.RegisterReq) error {
	if registerReq.Password != registerReq.Password2 {
		return fmt.Errorf("password is not the same")
	}

	hashedPassword, err := s.HashPassword(registerReq.Password)
	if err != nil {
		return err
	}

	user := models.User{
		Username: registerReq.Username,
		Password: hashedPassword,
	}
	err = s.db.Create(&user).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *UsersService) UpdateTokens(c *fiber.Ctx) error {
	cookieRefreshToken := c.Cookies("refresh_token")
	payload, err := s.VerifyToken(cookieRefreshToken)
	if err != nil {
		return err
	}

	dbRefreshToken := ""
	err = s.db.Model(&models.User{}).
		Where("id = ?", payload.ID).
		Select("refresh_token").
		First(&dbRefreshToken).
		Error
	if err != nil {
		return err
	}

	if cookieRefreshToken != dbRefreshToken {
		return fmt.Errorf("incorrect refresh token")
	}

	newAccessToken, err := s.GenerateToken(
		payload.ID,
		payload.Username,
		AccessTokenExpTime,
	)
	if err != nil {
		return err
	}

	newRefreshToken, err := s.GenerateToken(
		payload.ID,
		payload.Username,
		RefreshTokenExpTime,
	)
	if err != nil {
		return err
	}

	err = s.db.
		Model(&models.User{}).
		Where("id = ?", payload.ID).
		Update("refresh_token", newRefreshToken).
		Error
	if err != nil {
		return err
	}

	s.SetCookie(c, "access_token", newAccessToken, AccessTokenExpTime)
	s.SetCookie(c, "refresh_token", newRefreshToken, RefreshTokenExpTime)

	return nil
}

func (s *UsersService) VerifyTokenByCookie(
	c *fiber.Ctx,
	cookieName string,
) (*models.JwtPayload, error) {
	tokenString := c.Cookies(cookieName)
	if tokenString == "" {
		return nil, fmt.Errorf("token not found")
	}

	return s.VerifyToken(tokenString)
}

func (s *UsersService) VerifyToken(tokenString string) (*models.JwtPayload, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&models.JwtClaim{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(s.cfg.JwtSecret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.JwtClaim)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return &claims.Payload, nil
}

func (s *UsersService) ClearToken(c *fiber.Ctx) {
	s.SetCookie(c, "access_token", "deleted", -2*time.Hour)
	s.SetCookie(c, "refresh_token", "deleted", -2*time.Hour)
}

func (s *UsersService) SetCookie(
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

func (s *UsersService) GenerateToken(
	userId int,
	username string,
	expTime time.Duration,
) (string, error) {
	claims := models.JwtClaim{
		Payload: models.JwtPayload{
			ID:       userId,
			Username: username,
			Exp:      time.Now().Add(expTime).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	return token.SignedString([]byte(s.cfg.JwtSecret))
}

func (s *UsersService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *UsersService) CheckPassword(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (s *UsersService) IsAdmin(id int) (bool, error) {
	var isAdmin bool
	err := s.db.Model(&models.User{}).Where("id = ?", id).Select("is_admin").First(&isAdmin).Error
	if err != nil {
		return false, err
	}

	return isAdmin, nil
}
