package utils

import (
	"regexp"

	"gopkg.in/go-playground/validator.v9"
)

var (
	validate *validator.Validate
)

func InitValidators() {
	validate = validator.New()
	_ = validate.RegisterValidation("time", func(fl validator.FieldLevel) bool {
		compile, err := regexp.Compile("^[0-9]{4}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1]) (0[0-9]|1[0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9])$")
		if err != nil {
			panic("invalid validator")
		}
		return compile.MatchString(fl.Field().String())
	})
	_ = validate.RegisterValidation("email", func(fl validator.FieldLevel) bool {
		compile := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
		return compile.MatchString(fl.Field().String())
	})
	_ = validate.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		compile := regexp.MustCompile("^[a-z0-9-_.]+$")
		return compile.MatchString(fl.Field().String())
	})
}
