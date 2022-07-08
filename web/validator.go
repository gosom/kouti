package web

import (
	"github.com/go-playground/validator/v10"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("password", validatePassword)
}

func validatePassword(fl validator.FieldLevel) bool {
	const minEntropyBits = 50
	err := passwordvalidator.Validate(fl.Field().String(), minEntropyBits)
	return err == nil
}
