package utils

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			for _, e := range ve {
				return fmt.Errorf("field %s failed validation %s", e.Field(), e.Tag())
			}
		}
		var ivd *validator.InvalidValidationError
		if errors.As(err, &ivd) {
			return fmt.Errorf("invalid validation: %s", ivd)
		}
		return err
	}
	return nil
}

func FormatValidationErrors(err error) map[string]string {
	errType := err.(validator.ValidationErrors)

	for _, fieldError := range errType {
		switch fieldError.Tag() {
		case "required":
			return map[string]string{
				fieldError.Field(): fieldError.Field() + " is required",
			}
		case "email":
			return map[string]string{
				fieldError.Field(): fieldError.Field() + " must be a valid email",
			}
		case "min":
			return map[string]string{
				fieldError.Field(): fieldError.Field() + " must be at least " + fieldError.Param() + " characters",
			}
		case "max":
			return map[string]string{
				fieldError.Field(): fieldError.Field() + " must be at most " + fieldError.Param() + " characters",
			}
		}
	}

	return nil
}
