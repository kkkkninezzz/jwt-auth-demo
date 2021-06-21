package gossodemo

import (
	"fmt"

	"github.com/gofiber/fiber"
	"github.com/gomodule/redigo/redis"
)

func Boot() {
	c1, err := redis.Dial("tcp", "docker_redis:6379")
	if err != nil {
		panic(err)
	}
	rec1, err := c1.Do("Get", "gwyy")
	if err != nil {
		panic(err)
	}
	fmt.Println(rec1)

	defer c1.Close()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) {
		c.Send("Hello, World!")
	})

	app.Listen(3000)
}
