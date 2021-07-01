package authjwtdemo

import (
	"authjwtdemo/internal/app/authjwtdemo/config"
	"authjwtdemo/internal/app/authjwtdemo/model"
	"authjwtdemo/internal/app/authjwtdemo/router"
	"authjwtdemo/internal/pkg/database"
	"authjwtdemo/internal/pkg/redis"
	"log"

	"github.com/gofiber/fiber/v2"
)

var fiberApp *fiber.App

func Boot(configPath string) {
	config.Init(configPath)
	c := config.Config
	redis.Connect(c.RedisConfig.Host, c.RedisConfig.Port)
	// 初始化数据库
	database.Connect(c.MysqlConfig.Dsn, c.MysqlConfig.TablePrefix, c.MysqlConfig.SingularTable)
	model.Init()

	fiberApp = fiber.New()
	router.SetupRoutes(fiberApp)
	log.Fatalln(fiberApp.Listen(c.FiberAddr))
	log.Println("server is started!")

	defer redis.Shutdown()

}
