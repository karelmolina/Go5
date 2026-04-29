package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/karelmolina/play5/database"
	"github.com/karelmolina/play5/model"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func ValidateStruct(v interface{}) error {
	return validate.Struct(v)
}

func IsUsernameTaken(username string) bool {
	var count int64
	database.DB.Model(&model.User{}).Where("LOWER(username) = LOWER(?)", username).Count(&count)
	return count > 0
}

func ValidationErrorsToMap(err error) map[string]string {
	result := make(map[string]string)
	if verrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range verrs {
			field := strings.ToLower(e.Field())
			switch e.Tag() {
			case "required":
				result[field] = fmt.Sprintf("%s is required", field)
			case "min":
				result[field] = fmt.Sprintf("%s is too short", field)
			case "max":
				result[field] = fmt.Sprintf("%s is too long", field)
			default:
				result[field] = fmt.Sprintf("%s is invalid", field)
			}
		}
	}
	return result
}
