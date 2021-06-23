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
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	errors := validateStruct(out)
	if errors != nil {
		return ctx.JSON(errors)
	}

	return nil
}

func InitValidator() {
	validate = validator.New()
}
