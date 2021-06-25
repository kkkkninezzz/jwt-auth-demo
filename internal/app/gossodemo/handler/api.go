package handler

import "github.com/gofiber/fiber/v2"

func Hello(ctx *fiber.Ctx) error {
    return ctx.JSON(fiber.Map{"status": "success", "message": "Hello i'm ok!", "data": nil})
}
