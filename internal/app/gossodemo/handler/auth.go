package handler

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

type LoginInput struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func Login(ctx *fiber.Ctx) error {
	var input LoginInput
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	username := input.UserName
	password := input.Password
	if username != "kurisu9" || password != "123456" {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	mySigningKey := []byte("AllYourBase")
	t, err := token.SignedString(mySigningKey)
	if err != nil {
		log.Println(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.JSON(fiber.Map{"status": "success", "message": "Success login", "data": t})
}
