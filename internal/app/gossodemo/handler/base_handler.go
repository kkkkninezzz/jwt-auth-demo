package handler

import (
	"reflect"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	FailedField string `json:"failed_field"`
	Tag         string `json:"tag"`
	Value       string `json:"value"`
}

var validate *validator.Validate

func validateStruct(data interface{}) []*ErrorResponse {
	var errors []*ErrorResponse

	err := validate.Struct(data)
	if err != nil {
		types := reflect.TypeOf(data)
		// 如果data是指针，那么需要拿到原始的struct
		if types.Kind() == reflect.Ptr {
			types = types.Elem()
		}

		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse

			// 如果有json tag
			// 那么使用json tag的名称
			failedField := ""
			if structField, result := types.FieldByName(err.StructField()); result {
				failedField = structField.Tag.Get("json")
			}

			if failedField == "" {
				failedField = err.Field()
			}

			element.FailedField = failedField
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

// 从body中反序列化并且进行验证
func bodyParserAndValidate(out interface{}, ctx *fiber.Ctx) error {
	if err := ctx.BodyParser(out); err != nil {
		return BadRequestError(ctx, err)
	}

	errors := validateStruct(out)
	if errors != nil {
		return BadRequestError(ctx, errors)
	}

	return nil
}

func InitValidator() {
	validate = validator.New()
}

func Error(ctx *fiber.Ctx, status int, msg string, data interface{}) error {
	return ctx.Status(status).JSON(fiber.Map{
		"status":  "error",
		"message": msg,
		"data":    data,
	})
}

// 未验证的错误信息
func UnauthorizedError(ctx *fiber.Ctx, msg string, data interface{}) error {
	return Error(ctx, fiber.StatusUnauthorized, msg, data)
}

// 参数错误error信息
func BadRequestError(ctx *fiber.Ctx, data interface{}) error {
	return Error(ctx, fiber.StatusBadRequest, "Param is invaild", data)
}

// 服务器异常error信息
func InternalServerError(ctx *fiber.Ctx, msg string, data interface{}) error {
	return Error(ctx, fiber.StatusInternalServerError, msg, data)
}

// 成功的消息
func SuccessError(ctx *fiber.Ctx, msg string, data interface{}) error {
	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": msg,
		"data":    data,
	})
}
