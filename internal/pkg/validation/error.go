package validation

import "fmt"

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type FieldErrors []FieldError

func (fe FieldErrors) Error() string {
	var errStr string
	for _, fieldError := range fe {
		errStr += fmt.Sprintf("%s: %s\n", fieldError.Field, fieldError.Message)
	}
	return errStr
}
