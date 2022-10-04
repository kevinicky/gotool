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

	_ = validate.RegisterValidation("OVO", func(fl validator.FieldLevel) bool {
		return fl.Field().String() == "OVO"
	})

	_ = validate.RegisterValidation("DANA", func(fl validator.FieldLevel) bool {
		return fl.Field().String() == "DANA"
	})

	_ = validate.RegisterValidation("LINKAJA", func(fl validator.FieldLevel) bool {
		return fl.Field().String() == "LINKAJA"
	})

	_ = validate.RegisterValidation("SHOPEEPAY", func(fl validator.FieldLevel) bool {
		return fl.Field().String() == "SHOPEEPAY"
	})

	return
}
