package utils

import "github.com/go-playground/validator/v10"

func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		errors[err.Field()] = err.Tag()
	}
	return errors
}
