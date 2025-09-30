package handler

import (
	"fmt"
	"net/http"
	// "strings"

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
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
		}

		var errorMessage string
		for _, fieldError := range validationErrors {
			switch fieldError.Tag() {
			case "required":
				errorMessage = fmt.Sprintf("%s is required", fieldError.Field())
			case "email":
				errorMessage = fmt.Sprintf("%s must be a valid email address", fieldError.Field())
			case "min":
				errorMessage = fmt.Sprintf("%s must be at least %s characters long", fieldError.Field(), fieldError.Param())
			case "max":
				errorMessage = fmt.Sprintf("%s must be at most %s characters long", fieldError.Field(), fieldError.Param())
			default:
				errorMessage = fmt.Sprintf("validation failed on field %s with rule %s", fieldError.Field(), fieldError.Tag())
			}
			break
		}

		return echo.NewHTTPError(http.StatusBadRequest, ValidationErrorResponse{Error: errorMessage})

		// validationError := ValidationErrorResponse{
		// 	Error: fmt.Sprintf("validation failed: %s", err.Error()),
		// }
		// return echo.NewHTTPError(http.StatusBadRequest, validationError)
	}
	return nil
}

func NewValidator() *CustomValidator {
	return &CustomValidator{validator: validator.New()}
}