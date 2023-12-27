package users

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

const AccessTokenExpTime = 24 * time.Hour // TODO: change to 15 * time.Minute
const RefreshTokenExpTime = 30 * 24 * time.Hour

func (s *UsersService) Login(loginReq *models.LoginReq, c *fiber.Ctx) error {
	user := new(models.User)
	err := s.db.
		Where(&models.User{Username: loginReq.Username}).
		First(user).Error
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if !s.checkPassword(loginReq.Password, user.Password) {
		return fmt.Errorf("password incorrect")
	}

	accessToken, err := s.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		return err
	}

	refreshToken, err := s.GenerateRefreshToken(user.ID, user.Username)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(AccessTokenExpTime),
		HTTPOnly: true,
		Secure:   true,
		Path:     "/",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(RefreshTokenExpTime),
		HTTPOnly: true,
		Secure:   true,
		Path:     "/",
	})

	fmt.Printf(
		"accesstoken: %s | refreshtoken: %s\n",
		c.Cookies("access_token"),
		c.Cookies("refresh_token"),
	)

	return nil
}

func (s *UsersService) Register(registerReq *models.RegisterReq) error {
	if registerReq.Password != registerReq.Password2 {
		return fmt.Errorf("password is not the same")
	}

	hashedPassword, err := s.hashPassword(registerReq.Password)
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

// func (s *UsersService) RefreshToken(c *fiber.Ctx) error { // TODO
//
// }

// ---------------------------------------------------- //

func (s *UsersService) IsAuthenticated(c *fiber.Ctx) bool {
	accessToken := c.Cookies("access_token")
	if accessToken == "" {
		return false
	}

	_, err := s.VerifyToken(accessToken)
	return err == nil
}

func (s *UsersService) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func (s *UsersService) ClearToken(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "deleted",
		Expires:  time.Now().Add(-2 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		Path:     "/",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "deleted",
		Expires:  time.Now().Add(-2 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		Path:     "/",
	})
}

func (s *UsersService) GenerateAccessToken(userId uint, username string) (string, error) {
	return s.generateToken(userId, username, AccessTokenExpTime)
}

func (s *UsersService) GenerateRefreshToken(userId uint, username string) (string, error) {
	return s.generateToken(userId, username, RefreshTokenExpTime)
}

func (s *UsersService) generateToken(
	userId uint,
	username string,
	expTime time.Duration,
) (string, error) {
	claims := jwt.MapClaims{
		"id":       userId,
		"username": username,
		"exp":      time.Now().Add(expTime).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
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
