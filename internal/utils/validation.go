package utils

import "github.com/go-playground/validator/v10"

type (
	ErrorResponse struct {
		Error       bool
		FailedField string
		Tag         string
		Value       interface{}
	}

	Validator struct {
		validator *validator.Validate
	}

	GlobalErrorHandlerResp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
)

func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

func (v *Validator) Validate(data interface{}) []ErrorResponse {
	var validationErrors []ErrorResponse

	errs := v.validator.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem ErrorResponse
			elem.Error = true
			elem.FailedField = err.Field()
			elem.Tag = err.Tag()
			elem.Value = err.Value()

			validationErrors = append(validationErrors, elem)
		}
	}
	return validationErrors
}
