package validation

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestValidPhoneNumber(t *testing.T) {
	tests := []struct {
		name     string
		phone    string
		expected bool
	}{
		{name: "valid_phone_with_plus_and_digits", phone: "+62234567890", expected: true},
		{name: "valid_phone_with_long_digits", phone: "+62345678901234567890", expected: false},
		{name: "invalid_phone_missing_plus", phone: "6234567890", expected: false},
		{name: "invalid_phone_with_spaces", phone: "+623 456 7890", expected: false},
		{name: "invalid_phone_with_letters", phone: "+623ABC456", expected: false},
		{name: "empty_phone", phone: "", expected: false},
		{name: "invalid_phone_with_special_chars", phone: "+123-456-7890", expected: false},
	}

	v := validator.New()
	_ = v.RegisterValidation("phonenumber", ValidPhoneNumber)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Var(tt.phone, "phonenumber")
			got := err == nil

			if got != tt.expected {
				t.Errorf("ValidPhoneNumber(%q) = %v, want %v", tt.phone, got, tt.expected)
			}
		})
	}
}
