package middleware

import (
	"errors"
	"gossodemo/internal/app/gossodemo/model"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"golang.org/x/crypto/bcrypt"
)

// Protected protect routes
func Protected() func(*fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte("secret"),
		ErrorHandler: jwtError,
	})
}

func JWTAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil

	}
}

// 生成jwt的秘钥
func GenerateJwtSecret(salt string) (string, error) {
	staticSecret := "secret"

	bytes, err := bcrypt.GenerateFromPassword([]byte(salt+"."+staticSecret), 1)
	return string(bytes), err
}

// 生成jwt token
func GenerateJwtToken(userBase *model.UserBase, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = userBase.Username
	claims["user_id"] = userBase.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	return token.SignedString([]byte(secret))
}

func jwtFromHeader(c *fiber.Ctx, header string, authScheme string) (string, error) {
	auth := c.Get(header)
	l := len(authScheme)
	if len(auth) > l+1 && strings.EqualFold(auth[:l], authScheme) {
		return auth[l+1:], nil
	}
	return "", errors.New("Missing or malformed JWT")
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})

	} else {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
	}
}
