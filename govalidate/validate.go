package govalidate

import "github.com/go-playground/validator/v10"

type ValidateTools interface {
	CustomValidator() (validate *validator.Validate)
}

type validateTools struct{}

func NewValidateTools() ValidateTools {
	return &validateTools{}
}

func (v *validateTools) CustomValidator() (validate *validator.Validate) {
	validate = validator.New()

	_ = validate.RegisterValidation("android", func(fl validator.FieldLevel) bool {
		return fl.Field().String() == "android"
	})

	_ = validate.RegisterValidation("ios", func(fl validator.FieldLevel) bool {
		return fl.Field().String() == "ios"
	})

	return
}
