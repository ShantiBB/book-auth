package validate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"

	"auth/internal/http/lib/schema/response"
)

func Error(errs validator.ValidationErrors) response.Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.Tag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required", err.Field()))
		case "passwd":
			errMsgs = append(errMsgs, fmt.Sprintf("Password must have 8+ characters, "+
				"with a capital letter, number, and symbol"))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is invalid", err.Field()))
		}
	}

	return response.Response{
		Status: "error",
		Error:  strings.Join(errMsgs, ", "),
	}
}

func Password(fl validator.FieldLevel) bool {
	pass := fl.Field().String()
	if len([]rune(pass)) < 8 {
		return false
	}

	patterns := []string{`[a-z]`, `[A-Z]`, `\d`, `[!@#\$%\^&\*\?]`}
	for _, re := range patterns {
		if !regexp.MustCompile(re).MatchString(pass) {
			return false
		}
	}

	return true
}
