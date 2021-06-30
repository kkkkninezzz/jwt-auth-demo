package authjwtdemo

import (
	"authjwtdemo/internal/app/authjwtdemo/model"
	"authjwtdemo/internal/app/authjwtdemo/router"
	"authjwtdemo/internal/pkg/database"
	"authjwtdemo/internal/pkg/redis"
	"log"

	"github.com/gofiber/fiber/v2"
)

var fiberApp *fiber.App

func Boot() {
	redis.Connect("127.0.0.1", 6379)
	// 初始化数据库
	database.Connect("root:123456@tcp(127.0.0.1:3306)/demo-orm?charset=utf8mb4&parseTime=True&loc=Local")
	model.Init()

	fiberApp = fiber.New()
	router.SetupRoutes(fiberApp)
	log.Fatalln(fiberApp.Listen(":3000"))
	log.Println("server is started!")

	defer redis.Shutdown()

}
