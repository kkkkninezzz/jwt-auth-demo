package router

import (
	"authjwtdemo/internal/app/authjwtdemo/handler"
	"authjwtdemo/internal/app/authjwtdemo/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	handler.InitValidator()

	api := app.Group("/api", logger.New())
	api.Get("/", handler.Hello)

	// auth
	auth := api.Group("/auth")
	auth.Post("/login", handler.Login)
	auth.Post("/register", handler.Register)

	// Products
	product := api.Group("/product")
	product.Get("/", middleware.Protected(), handler.GetAllProducts)
	product.Post("/", middleware.Protected(), handler.CreateProduct)
}
