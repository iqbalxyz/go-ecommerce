package utils

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

func FormatValidationErrors(err error) map[string]string {
	var errs validator.ValidationErrors
	if !errors.As(err, &errs) {
		return map[string]string{
			"_error": "unknown validation error",
		}
	}

	result := make(map[string]string)
	for _, fieldError := range errs {
		result[fieldError.Field()] = messageForTag(fieldError)
	}
	return result
}

func messageForTag(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return e.Field() + " is required"
	case "email":
		return e.Field() + " must be a valid email"
	case "min":
		return e.Field() + " must be at least " + e.Param() + " characters"
	case "max":
		return e.Field() + " must be at most " + e.Param() + " characters"
	default:
		return e.Field() + " is invalid"
	}
}
