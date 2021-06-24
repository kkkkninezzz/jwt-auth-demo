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

		pMapClaims, ok := t.Claims.(*jwt.MapClaims)
		if !ok {
			return nil, errors.New("Not support Claims")
		}

		mapClaims := *pMapClaims
		userId, keyExists := mapClaims["user_id"]
		if !keyExists {
			return nil, errors.New("User_id is not in Claims")
		}

		var uid uint
		switch v := userId.(type) {
		case float64:
			uid = uint(v)
		default:
			return nil, errors.New("User_id type is not uint")
		}

		salt := redis.Template.Get(rediskey.FormatSaltRedisKey(uid))
		if salt == "" {
			return nil, errors.New("Not found salt")
		}

		secret, err := middleware.GenerateJwtSecret(salt)
		if err != nil {
			return nil, err
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
