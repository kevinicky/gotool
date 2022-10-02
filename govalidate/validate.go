package govalidate

import "github.com/go-playground/validator/v10"

type ValidatorUtil interface {
	CustomValidator() (validate *validator.Validate)
}

type validatorUtil struct{}

func NewValidate() ValidatorUtil {
	return &validatorUtil{}
}

func (v *validatorUtil) CustomValidator() (validate *validator.Validate) {
	validate = validator.New()

	_ = validate.RegisterValidation("android", func(fl validator.FieldLevel) bool {
		return fl.Field().String() == "android"
	})

	_ = validate.RegisterValidation("ios", func(fl validator.FieldLevel) bool {
		return fl.Field().String() == "ios"
	})

	return
}
