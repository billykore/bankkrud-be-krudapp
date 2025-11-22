package validation

import (
	"reflect"
	"strings"
)

// tagMessages maps validation tags to corresponding error message templates.
var tagMessages = map[string]string{
	"required":    "required",
	"email":       "not a valid email",
	"len":         "length must be %s",
	"min":         "minimum length must be %s",
	"number":      "must be a number",
	"gte":         "must be greater than or equal to %s",
	"lte":         "must be less than or equal to %s",
	"phonenumber": "not a valid phone number",
	"only":        "must contain only: %s",
	"oneof":       "must be one of: %s",
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
