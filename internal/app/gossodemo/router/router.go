package router

import (
	"gossodemo/internal/app/gossodemo/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api", logger.New())
	api.Get("/", handler.Hello)

	// auth
	auth := app.Group("/auth")
	auth.Post("/login", handler.Login)
}
