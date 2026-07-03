package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	v *validator.Validate
}

func New() Validator {
	return Validator{v: validator.New()}
}

func (val Validator) ValidateStruct(value any) error {
	if err := val.v.Struct(value); err != nil {
		return err
	}
	return nil
}

func FieldErrors(err error) map[string]string {
	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return map[string]string{"_": err.Error()}
	}
	out := make(map[string]string, len(validationErrs))
	for _, fe := range validationErrs {
		field := strings.ToLower(fe.Field())
		out[field] = fmt.Sprintf("failed on '%s' tag", fe.Tag())
	}
	return out
}
