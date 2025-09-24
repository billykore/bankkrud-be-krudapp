package validation

import "fmt"

type FieldError struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

type FieldErrors []FieldError

func (fe FieldErrors) Error() string {
	var errStr string
	for _, fieldError := range fe {
		errStr += fmt.Sprintf("%s: %s\n", fieldError.Name, fieldError.Message)
	}
	return errStr
}
