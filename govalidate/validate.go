package govalidate

import "github.com/go-playground/validator/v10"

type ValidateTools interface {
	CustomValidator() (validate *validator.Validate)
	CustomValidationError(error error) map[string]interface{}
}

type validateTools struct{}

// NewValidateTools create object for accessing validation tool function
//
// Return validate tool interface for accessing function
func NewValidateTools() ValidateTools {
	return &validateTools{}
}

// CustomValidator registers custom tag for validator/v10
//
// Return validate object
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

// CustomValidationError validates error from struct tag that earlier
// registered.
//
// Return message error
func (v *validateTools) CustomValidationError(error error) map[string]interface{} {
	message := map[string]interface{}{}
	if castedObject, ok := error.(validator.ValidationErrors); ok {
		errObj := castedObject[0]

		switch errObj.Tag() {
		case "required":
			message = map[string]interface{}{"error": errObj.Field() + " is required"}
		case "android|ios":
			message = map[string]interface{}{"error": errObj.Field() + " must android, ios"}
		case "DANA|LINKAJA|OVO|SHOPEEPAY":
			message = map[string]interface{}{"error": errObj.Field() + " must DANA, LINKAJA, OVO, or SHOPEEPAY"}
		case "email":
			message = map[string]interface{}{"error": errObj.Field() + " is not valid email format"}
		case "gte":
			message = map[string]interface{}{"error": errObj.Field() + " value must be greater equal than " + errObj.Param()}
		case "gt":
			message = map[string]interface{}{"error": errObj.Field() + " value must be greater than " + errObj.Param()}
		case "lte":
			message = map[string]interface{}{"error": errObj.Field() + " value must be less equal than " + errObj.Param()}
		case "lt":
			message = map[string]interface{}{"error": errObj.Field() + " value must be less than " + errObj.Param()}
		default:
			message = map[string]interface{}{"error": "invalid input for " + errObj.Field()}
		}
	}

	return message
}
