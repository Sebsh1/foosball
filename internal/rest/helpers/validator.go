package helpers

import (
	validator "github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()

	return &Validator{
		validator: v,
	}
}

func (cv *Validator) Validate(i interface{}) error {
	return errors.Wrap(cv.validator.Struct(i), "validation failed")
}
