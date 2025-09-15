package config

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func NewValidator(viper *viper.Viper) *validator.Validate {
	v := validator.New()

	v.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		re := regexp.MustCompile(`^\+?[0-9]{8,15}$`)
		return re.MatchString(phone)
	})
	return v
}
