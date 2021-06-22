package gossodemo

import (
	"gossodemo/internal/pkg/redis"

	"github.com/gofiber/fiber/v2"
)

var redisTemplate redis.RedisTemplate
var fiberApp *fiber.App

func Boot() {
	redisTemplate = redis.NewRedisTemplate("127.0.0.1", 6379)

	fiberApp = fiber.New()
	fiberApp.Listen(":3000")

}

func Stop() {
	redisTemplate.Close()
	fiberApp.Shutdown()
}
