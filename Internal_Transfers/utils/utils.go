package utils

import (
	"internal-transfers/dtos"

	"github.com/go-playground/validator/v10"
)

func ValidateStruct(body interface{}) []*dtos.ErrorResponse {
	var errors []*dtos.ErrorResponse
	validate := validator.New()
	err := validate.Struct(body)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element dtos.ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}