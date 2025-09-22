package validation

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Validator is a struct that provides validation functionality for request data.
type Validator struct {
	v *validator.Validate
}

// New initializes and returns a new instance of Validator.
func New() *Validator {
	vv := &Validator{
		v: validator.New(),
	}
	if err := vv.registerCustomValidation(); err != nil {
		panic(err)
	}
	vv.JSONTagFunc()
	return vv
}

// Validate checks the validity of the input data based on struct tags and returns an error if any.
func (v *Validator) Validate(req any) error {
	err := v.v.Struct(req)
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return nil
	}
	return joinValidationErrors(ve)
}

// joinValidationErrors converts a slice of FieldError into a concatenated error.
func joinValidationErrors(validationErrors validator.ValidationErrors) error {
	var fieldErrors FieldErrors
	for _, fieldError := range validationErrors {
		fieldErrors = append(fieldErrors, FieldError{
			Field:   fieldError.Field(),
			Message: errMessage(fieldError),
		})
	}
	return fieldErrors
}

// errMessage formats and returns error messages based on the field validation error type.
func errMessage(fe validator.FieldError) string {
	if messageTemplate, exists := tagMessages[fe.Tag()]; exists {
		return formatMessage(messageTemplate, fe.Field(), fe.Param())
	}
	return formatMessage("%s format is invalid", fe.Field())
}

// formatMessage is a helper function to format messages with parameters.
func formatMessage(template, field string, param ...string) string {
	if len(param) > 0 && param[0] != "" {
		return fmt.Sprintf(template, field, param[0])
	}
	return fmt.Sprintf(template, field)
}
