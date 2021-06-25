package rediskey

import "fmt"

func FormatJwtSaltRedisKey(id uint) string {
	return fmt.Sprintf("user_jwt_salt:%d", id)
}
