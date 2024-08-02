package responce

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Responce struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOk    = "OK"
	StatusError = "Error"
)

func OK() Responce {
	return Responce{
		Status: StatusOk,
	}
}

func Error(msg string) Responce {
	return Responce{
		Status: StatusError,
		Error:  msg,
	}
}

func ValidatorError(errs validator.ValidationErrors) Responce {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("filed %s is a required filed", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("filed %s is a valid URL", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("filed %s is a valid", err.Field()))
		}
	}

	return Responce{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
