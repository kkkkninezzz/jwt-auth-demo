package handler

import (
    "gossodemo/internal/app/gossodemo/model"
    "gossodemo/internal/pkg/database"

    "github.com/gofiber/fiber/v2"
)

// GetAllProducts query all products
func GetAllProducts(c *fiber.Ctx) error {
    db := database.DB
    var products []model.Product
    db.Find(&products)
    return c.JSON(fiber.Map{"status": "success", "message": "All products", "data": products})
}

// CreateProduct new product
func CreateProduct(c *fiber.Ctx) error {
    db := database.DB
    product := new(model.Product)
    if err := c.BodyParser(product); err != nil {
        return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't create product", "data": err})
    }
    db.Create(&product)
    return c.JSON(fiber.Map{"status": "success", "message": "Created product", "data": product})
}
