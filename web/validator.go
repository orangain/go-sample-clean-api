package web

import (
	"github.com/labstack/echo"
	validator "gopkg.in/go-playground/validator.v9"
)

func NewEchoCustomValidator() echo.Validator {
	return &CustomValidator{validator: validator.New()}
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
