package gossodemo

import (
	"gossodemo/internal/app/gossodemo/router"
	"gossodemo/internal/pkg/redis"
	"log"

	"github.com/gofiber/fiber/v2"
)

var fiberApp *fiber.App

func Boot() {
	redis.Connect("127.0.0.1", 6379)
	// TODO 初始化数据库

	fiberApp = fiber.New()
	router.SetupRoutes(fiberApp)
	log.Fatalln(fiberApp.Listen(":3000"))

	defer redis.Shutdown()

}
