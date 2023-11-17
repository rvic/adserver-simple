package utils

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// NewValidator func for create a new validator for model fields.
func NewValidator() *validator.Validate {
	validate := validator.New()

	// Custom validation for uuid.UUID fields.
	_ = validate.RegisterValidation("uuid", func(fl validator.FieldLevel) bool {
		field := fl.Field().String()
		if _, err := uuid.Parse(field); err != nil {
			return false // if there is an error, validation should return false
		}
		return true // if no error, validation should return true
	})

	return validate
}

func ValidatorErrors(err error) map[string]string {
	var errs validator.ValidationErrors
	if ok := errors.As(err, &errs); !ok {
		return nil
	}

	fields := make(map[string]string, len(errs))
	for _, err := range errs {
		fields[err.Field()] = err.Error()
	}

	return fields
}
