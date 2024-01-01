package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidateStruct(s any, val *validator.Validate) error {
	if val == nil {
		val = validator.New()
	}
	if err := val.Struct(s); err != nil {
		var errors []string
		for _, e := range err.(validator.ValidationErrors) {
			newErr := fmt.Sprintf("%s failed on the %s tag with value %s", e.Field(), e.Tag(), e.Value())
			errors = append(errors, newErr)
		}
		return fmt.Errorf(strings.Join(errors, ", "))
	}
	return nil
}

func ValidateStructAll(s any, validate *validator.Validate) []error {
	if validate == nil {
		validate = validator.New()
	}
	var errors []error
	if err := validate.Struct(s); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			newErr := fmt.Sprintf("%s failed on the %s tag with value %s", e.Field(), e.Tag(), e.Value())
			errors = append(errors, fmt.Errorf(newErr))
		}
	}
	return errors
}
