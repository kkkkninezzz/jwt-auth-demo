package middleware

import (
	"crypto/md5"
	"errors"
	"fmt"
	"gossodemo/internal/app/gossodemo/def/rediskey"
	"gossodemo/internal/app/gossodemo/model"
	"gossodemo/internal/pkg/redis"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

const TokenKey string = "user"
const JWTAuthScheme string = "Bearer"

// Protected protect routes
func Protected() func(*fiber.Ctx) error {
	/*
		return jwtware.New(jwtware.Config{
			SigningKey:   []byte("secret"),
			ErrorHandler: jwtError,
		})
	*/
	return JWTAuthMiddleware()
}

func JWTAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth, err := jwtFromHeader(c, fiber.HeaderAuthorization, JWTAuthScheme)
		if err != nil {
			return jwtError(c, err)
		}

		token, err := jwt.ParseWithClaims(auth, &jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
			pMapClaims, ok := t.Claims.(*jwt.MapClaims)
			if !ok {
				return nil, errors.New("not support Claims")
			}

			mapClaims := *pMapClaims
			userId, keyExists := mapClaims["user_id"]
			if !keyExists {
				return nil, errors.New("user_id is not in Claims")
			}

			var uid uint
			switch v := userId.(type) {
			case float64:
				uid = uint(v)
			default:
				return nil, errors.New("user_id type is not uint")
			}

			salt := redis.Template.Get(rediskey.FormatSaltRedisKey(uid))
			if salt == "" {
				return nil, errors.New("not found salt")
			}

			secret := GenerateJwtSecret(salt)
			if secret == "" {
				return nil, errors.New("secret generate failed")
			}
			return []byte(secret), nil
		})

		if err == nil && token.Valid {
			// Store user information from token into context.
			c.Locals(TokenKey, token)
			return c.Next()
		}

		return jwtError(c, err)
	}
}

// 生成jwt的秘钥
func GenerateJwtSecret(salt string) string {
	staticSecret := "secret"
	return fmt.Sprintf("%x", md5.Sum([]byte(salt+"."+staticSecret)))
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
