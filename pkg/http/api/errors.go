package api

import (
	"fmt"

	"github.com/go-playground/validator"
)

func UnexpectedError() Response {
	return Err("An unexpected error occured.")
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errDetails []ErrDetail

	for _, err := range errs {
		field := err.Field()

		var msg string

		switch err.ActualTag() {
		case "required":
			msg = fmt.Sprintf("field %s is required.", field)
		case "min":
			msg = fmt.Sprintf("field %s min length must be %s.", field, err.Param())
		case "max":
			msg = fmt.Sprintf("field %s max length must be %s.", field, err.Param())
		case "alphanum":
			msg = fmt.Sprintf("field %s must contain both letters and numbers.", field)
		case "containsany":
			msg = fmt.Sprintf("field %s must contain on of the following characters: %s.", field, err.Param())
		case "eqfield":
			msg = fmt.Sprintf("field %s is not equal to %s field.", field, err.Param())
		case "passwordpattern":
			msg = "field password must contain at least one letter, one number, and one special character."
		default:
			msg = fmt.Sprintf("field %s is invalid", field)
		}

		errDetails = append(errDetails, ErrDetail{Field: field, Info: msg})
	}

	return ErrD("validation error", errDetails)
}
