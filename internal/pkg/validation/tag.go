package validation

import (
	"reflect"
	"strings"
)

// tagMessages maps validation tags to corresponding error message templates.
var tagMessages = map[string]string{
	"required": "%s is required",
	"email":    "%s is not a valid email",
	"len":      "%s length must be %s",
	"min":      "%s minimum length must be %s",
	"number":   "%s must be a number",
	"gte":      "%s must be greater than or equal to %s",
	"lte":      "%s must be less than or equal to %s",
}

func (v *Validator) JSONTagFunc() {
	v.v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" { // Handle ignored JSON fields
			return ""
		}
		return name
	})
}
