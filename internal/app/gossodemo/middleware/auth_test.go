package middleware_test

import (
	"errors"
	"gossodemo/internal/app/gossodemo/def/rediskey"
	"gossodemo/internal/app/gossodemo/middleware"
	"gossodemo/internal/pkg/redis"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

func TestJwt(t *testing.T) {
	redis.Connect("127.0.0.1", 6379)
	auth := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjQ3ODgxNTAsInVzZXJfaWQiOjEsInVzZXJuYW1lIjoia3VyaXN1OSJ9.39Nxbq5R8TMqIWnb98ch-ZYi48-lfT9yU6Y5_y2teO0"

	token, err := jwt.ParseWithClaims(auth, &jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {

		userClaims, ok := t.Claims.(*middleware.UserClaims)
		if !ok {
			return nil, errors.New("not support Claims")
		}

		userId := userClaims.UserInfo.UserId
		if userId <= 0 {
			return nil, errors.New("missing or malformed JWT")
		}

		salt := redis.Template.Get(rediskey.FormatJwtSaltRedisKey(userId))
		if salt == "" {
			return nil, errors.New("Not found salt")
		}

		secret := middleware.GenerateJwtSecret(salt)
		if secret == "" {
			return nil, errors.New("Secret generate failed")
		}
		return []byte(secret), nil
	})

	if err != nil {
		t.Error(err)
	}

	t.Log(token)
}
