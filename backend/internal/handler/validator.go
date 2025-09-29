package handler

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

type ValidationErrorResponse struct {
	Error string `json:"error"`
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		validationError := ValidationErrorResponse{
			Error: fmt.Sprintf("validation failed: %s", err.Error()),
		}
		return echo.NewHTTPError(http.StatusBadRequest, validationError)
	}
	return nil
}

func NewValidator() *CustomValidator {
	return &CustomValidator{validator: validator.New()}
}