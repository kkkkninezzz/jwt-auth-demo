package handler

import (
	"gossodemo/internal/app/gossodemo/def/rediskey"
	"gossodemo/internal/app/gossodemo/middleware"
	"gossodemo/internal/app/gossodemo/model"
	"gossodemo/internal/pkg/database"
	"gossodemo/internal/pkg/random/randstr"
	"gossodemo/internal/pkg/redis"
	"log"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	UserName string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=3,max=20"`
}

func checkPasswordHash(userBase *model.UserBase, originPassword string) bool {
	// CompareHashAndPassword这方法是真滴慢，估计得考虑降低cost
	err := bcrypt.CompareHashAndPassword([]byte(userBase.Password), []byte(addSaltToPassword(originPassword, userBase.Salt)))
	return err == nil
}

func getUserByUsername(u string) (*model.UserBase, error) {
	db := database.DB
	var userBase model.UserBase
	result := db.Where(&model.UserBase{Username: u}).Find(&userBase)
	if result.RowsAffected <= 0 {
		return nil, nil
	}

	if err := result.Error; err != nil {
		if database.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	return &userBase, nil
}

func Login(ctx *fiber.Ctx) error {
	var input LoginInput
	if err := bodyParserAndValidate(&input, ctx); err != nil {
		return err
	}

	username := input.UserName
	password := input.Password

	userBase, err := getUserByUsername(username)
	if err != nil {
		return UnauthorizedError(ctx, "Error on username", err)
	}

	if userBase == nil {
		return UnauthorizedError(ctx, "User not found", username)
	}

	if !checkPasswordHash(userBase, password) {
		return UnauthorizedError(ctx, "Invalid password", nil)
	}

	secret := middleware.GenerateJwtSecret(userBase.Salt)
	if secret == "" {
		return UnauthorizedError(ctx, "Generate secret failed", nil)
	}

	t, err := middleware.GenerateJwtToken(userBase, secret)
	if err != nil {
		log.Println(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	// 写入redis
	// TODO 是否要进行过时的处理
	redis.Template.Set(rediskey.FormatSaltRedisKey(userBase.ID), userBase.Salt)
	return SuccessError(ctx, "Success login", t)
}

type RegisterInput struct {
	UserName string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=3,max=20"`
}

func Register(ctx *fiber.Ctx) error {
	var input RegisterInput
	if err := bodyParserAndValidate(&input, ctx); err != nil {
		return err
	}

	username := input.UserName
	password := input.Password

	userBase := new(model.UserBase)
	userBase.Username = username

	// 生成salt
	salt, err := generateSalt(username)
	if err != nil {
		return InternalServerError(ctx, "Couldn't generate salt", err)
	}
	userBase.Salt = salt

	// 加密
	encryptPassword, err := hashPassword(password, userBase.Salt)
	if err != nil {
		return InternalServerError(ctx, "Couldn't hash password", err)
	}
	userBase.Password = encryptPassword

	db := database.DB
	if err := db.Create(&userBase).Error; err != nil {
		return InternalServerError(ctx, "Couldn't create user", err)
	}

	return SuccessError(ctx, "Register Success", input)
}

func generateSalt(username string) (string, error) {
	// TODO 读取配置
	str := randstr.RandomAscii(20)

	bytes, err := bcrypt.GenerateFromPassword([]byte(username+"."+str), 14)
	return string(bytes), err
}

func addSaltToPassword(password string, salt string) string {
	return password + "." + salt
}

func hashPassword(password string, salt string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(addSaltToPassword(password, salt)), 14)
	return string(bytes), err
}
