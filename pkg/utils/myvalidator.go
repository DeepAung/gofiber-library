package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type MyValidator struct {
	validator *validator.Validate
}

func NewMyValidator() *MyValidator {
	return &MyValidator{
		validator: validator.New(),
	}
}

func (m *MyValidator) Validate(mystruct interface{}) error {
	err := m.validator.Struct(mystruct)
	if err == nil {
		return nil
	}

	msg := ""
	for _, e := range err.(validator.ValidationErrors) {
		msg += fmt.Sprintf("[%s]: '%v' | Needs to implement '%s'\n", e.Field(), e.Value(), e.Tag())
	}

	if msg == "" {
		return nil
	}

	return fmt.Errorf(msg)
}
