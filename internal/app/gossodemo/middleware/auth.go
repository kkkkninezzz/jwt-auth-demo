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

const UserInfoKey string = "user"
const JWTAuthScheme string = "Bearer"

type UserSimpleInfo struct {
	Username string `json:"username"`
	UserId   uint   `json:"user_id"`
}

type UserClaims struct {
	UserInfo UserSimpleInfo `json:"user_info"`
	jwt.StandardClaims
}

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

		token, err := jwt.ParseWithClaims(auth, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
			userClaims, ok := t.Claims.(*UserClaims)
			if !ok {
				return nil, errors.New("not support Claims")
			}

			userId := userClaims.UserInfo.UserId
			if userId <= 0 {
				return nil, errors.New("missing or malformed JWT")
			}

			salt := redis.Template.Get(rediskey.FormatSaltRedisKey(userId))
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
			userClaims := token.Claims.(*UserClaims)
			c.Locals(UserInfoKey, userClaims.UserInfo)
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
	claims := UserClaims{
		UserSimpleInfo{
			Username: userBase.Username,
			UserId:   userBase.ID,
		},
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			Issuer:    "go-sso-demo",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func jwtFromHeader(c *fiber.Ctx, header string, authScheme string) (string, error) {
	auth := c.Get(header)
	l := len(authScheme)
	if len(auth) > l+1 && strings.EqualFold(auth[:l], authScheme) {
		return auth[l+1:], nil
	}
	return "", errors.New("missing or malformed JWT")
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "missing or malformed JWT" {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})

	} else {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
	}
}

// 获取从jwt中解析得到的用户信息
func LocalUserInfo(c *fiber.Ctx) UserSimpleInfo {
	return c.Locals(UserInfoKey).(UserSimpleInfo)
}
