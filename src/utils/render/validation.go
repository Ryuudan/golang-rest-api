package render

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

func Validator() *validator.Validate {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}
		return name
	})

	// validate.RegisterValidation("dateComparison", dateComparison)
	return validate
}

func CustomValidationError(w http.ResponseWriter, r *http.Request, details []ValidationErrorDetails) {

	errorResponse := ValidationErrorResponse{
		Error: validationError{
			Code:      http.StatusUnprocessableEntity,
			Type:      "validation_error",
			Path:      r.URL.Path,
			TimeStamp: time.Now(),
			Details:   details,
		},
	}

	response, err := json.Marshal(errorResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	if _, err := w.Write(response); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

// ideal for form validation, or form errorrs
// like setError in react-hooks-form
func ValidationError(w http.ResponseWriter, r *http.Request, err error) {
	var details []ValidationErrorDetails

	for _, err := range err.(validator.ValidationErrors) {

		field := strings.ToLower(err.Field())
		errorType := validationErrorMessage(err)
		details = append(details, ValidationErrorDetails{Field: field, Message: errorType})
	}

	errorResponse := ValidationErrorResponse{
		Error: validationError{
			Code:      http.StatusUnprocessableEntity,
			Type:      "validation_error",
			TimeStamp: time.Now(),
			Path:      r.URL.Path,
			Details:   details,
		},
	}

	response, err := json.Marshal(errorResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	if _, err := w.Write(response); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func validationErrorMessage(err validator.FieldError) string {
	field := err.Field()
	println(err.Field())
	println(err.StructNamespace())

	switch tag := err.Tag(); tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s should be at least %s characters", field, err.Param())
	case "max":
		return fmt.Sprintf("%s should be at most %s characters", field, err.Param())
	case "lte":
		return fmt.Sprintf("%s should be less than or equal to %s", field, err.Param())
	case "gte":
		return fmt.Sprintf("%s should be greater than or equal to %s", field, err.Param())
	case "email":
		return fmt.Sprintf("%s is not a valid email address", field)
	case "url":
		return fmt.Sprintf("%s is not a valid URL", field)
	case "oneof":
		return fmt.Sprintf("%s should be one of [%s]", field, err.Param())
	case "alpha":
		return fmt.Sprintf("%s should only contain alphabetic characters", field)
	case "alphanum":
		return fmt.Sprintf("%s should only contain alphanumeric characters", field)
	case "numeric":
		return fmt.Sprintf("%s should be a numeric value", field)
	case "number":
		return fmt.Sprintf("%s should be a number value", field)
	case "eq":
		return fmt.Sprintf("%s should be equal to %s", field, err.Param())
	case "eq_ignore_case":
		return fmt.Sprintf("%s should be equal to %s (ignoring case)", field, err.Param())
	case "gt":
		return fmt.Sprintf("%s should be greater than %s", field, err.Param())
	case "lt":
		return fmt.Sprintf("%s should be less than %s", field, err.Param())
	case "ne":
		return fmt.Sprintf("%s should not be equal to %s", field, err.Param())
	case "ne_ignore_case":
		return fmt.Sprintf("%s should not be equal to %s (ignoring case)", field, err.Param())
	case "eqcsfield":
		return fmt.Sprintf("%s must be the same as %s", field, err.Param())
	case "eqfield":
		return fmt.Sprintf("%s must be the same as %s", field, err.Param())
	case "fieldcontains":
		return fmt.Sprintf("%s should contain %s", field, err.Param())
	case "fieldexcludes":
		return fmt.Sprintf("%s should not contain %s", field, err.Param())
	case "gtcsfield":
		return fmt.Sprintf("%s must be greater than %s", field, err.Param())
	case "gtecsfield":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, err.Param())
	case "gtefield":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, err.Param())
	case "gtfield":
		return fmt.Sprintf("%s must be greater than %s", field, err.Param())
	case "ltcsfield":
		return fmt.Sprintf("%s must be less than %s", field, err.Param())
	case "ltecsfield":
		return fmt.Sprintf("%s must be less than or equal to %s", field, err.Param())
	case "ltefield":
		return fmt.Sprintf("%s must be less than or equal to %s", field, err.Param())
	case "ltfield":
		return fmt.Sprintf("%s must be less than %s", field, err.Param())
	case "necsfield":
		return fmt.Sprintf("%s must be different from %s", field, err.Param())
	case "nefield":
		return fmt.Sprintf("%s must be different from %s", field, err.Param())
	case "e164":
		return fmt.Sprintf("%s must be a valid phone number", field)
	default:
		return fmt.Sprintf("Validation failed for %s field", field)
	}
}

type ValidationErrorDetails struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type validationError struct {
	Code      int                      `json:"code"`
	Type      string                   `json:"type"`
	Path      string                   `json:"path"`
	TimeStamp time.Time                `json:"time_stamp"`
	Details   []ValidationErrorDetails `json:"details"`
}

type ValidationErrorResponse struct {
	Error validationError `json:"error"`
}
