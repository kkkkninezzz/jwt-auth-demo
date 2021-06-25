package rediskey

import "fmt"

func FormatSaltRedisKey(id uint) string {
    return fmt.Sprintf("user_salt:%d", id)
}
