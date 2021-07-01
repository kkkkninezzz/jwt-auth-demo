package middleware

import (
	"authjwtdemo/internal/app/authjwtdemo/config"
	"authjwtdemo/internal/app/authjwtdemo/def/rediskey"
	"authjwtdemo/internal/app/authjwtdemo/model"
	"authjwtdemo/internal/pkg/redis"
	"authjwtdemo/internal/pkg/timeutil"
	"crypto/md5"
	"errors"
	"fmt"
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
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
		}

		token, err := jwt.ParseWithClaims(auth, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
			userClaims, ok := t.Claims.(*UserClaims)
			if !ok {
				return nil, errors.New("not support Claims")
			}

			userId := userClaims.UserInfo.UserId
			if userId <= 0 {
				return nil, errors.New("invalid user id ")
			}

			salt := redis.Template.Get(rediskey.FormatJwtSaltRedisKey(userId))
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
			c.Locals(UserInfoKey, &userClaims.UserInfo)
			return c.Next()
		}

		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
	}
}

// 生成用于jwt的密匙 salt
func GenerateJwtSecretSalt(userSalt string) string {
	return fmt.Sprintf("%s.%d", userSalt, timeutil.CurrentTimeMillis())
}

// 生成jwt的秘钥
func GenerateJwtSecret(salt string) string {
	staticSecret := config.Config.JwtConfig.PrivateSecret
	return fmt.Sprintf("%x", md5.Sum([]byte(salt+"."+staticSecret)))
}

// 生成jwt token
func GenerateJwtToken(userBase *model.UserBase, secret string, expiration time.Duration) (string, error) {
	claims := UserClaims{
		UserSimpleInfo{
			Username: userBase.Username,
			UserId:   userBase.ID,
		},
		jwt.StandardClaims{
			ExpiresAt: timeutil.NextTimeSeconds(expiration),
			Issuer:    config.Config.JwtConfig.Issuer,
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

// 获取从jwt中解析得到的用户信息
func LocalUserInfo(c *fiber.Ctx) UserSimpleInfo {
	return c.Locals(UserInfoKey).(UserSimpleInfo)
}
