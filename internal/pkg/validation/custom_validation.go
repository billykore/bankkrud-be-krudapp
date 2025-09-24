package validation

import (
	"regexp"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
)

// validationRegistry maps validation tags to their corresponding validation functions.
// It provides a centralized registry for custom field validators.
//
// Example usage:
//
//	registry := ValidationRegistry{
//	    "phonenumber": ValidatePhoneNumber,
//	}
type validationRegistry map[string]func(fl validator.FieldLevel) bool

// customValidations is a registry of custom validation functions.
//
// Add custom validation functions here.
var customValidations = validationRegistry{
	"phonenumber": ValidPhoneNumber,
	"only":        Only,
}

func ValidPhoneNumber(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	exp := regexp.MustCompile(`^\+62[0-9]{8,12}$`)
	return exp.MatchString(phone)
}

func (v *Validator) registerCustomValidation() error {
	for tag, fn := range customValidations {
		if err := v.v.RegisterValidation(tag, fn); err != nil {
			return err
		}
	}
	return nil
}

func Only(fl validator.FieldLevel) bool {
	values := strings.Split(fl.Field().String(), ",")
	params := strings.Split(fl.Param(), " ")
	for _, value := range values {
		if !slices.Contains(params, value) {
			return false
		}
	}
	return true
}
