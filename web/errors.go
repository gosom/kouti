package web

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type ErrHTTP struct {
	StatusCode int      `json:"-"`
	Message    string   `json:"message,omitempty"`
	Errors     []string `json:"errors,omitempty"`
}

func (e ErrHTTP) Error() string {
	if e.StatusCode > 0 {
		msg := fmt.Sprintf("%d: %s", e.StatusCode, e.Message)
		return msg
	}
	return e.Message
}

func NewBadRequestError(msg string, args ...string) ErrHTTP {
	ans := ErrHTTP{
		StatusCode: http.StatusBadRequest,
	}
	if len(msg) > 0 {
		arguments := make([]any, len(args))
		for i := range args {
			arguments[i] = args[i]
		}
		ans.Message = fmt.Sprintf(msg, arguments...)
	} else {
		ans.Message = http.StatusText(http.StatusBadRequest)
	}
	return ans
}

func NewValidationError(err error) ErrHTTP {
	ans := ErrHTTP{
		StatusCode: http.StatusBadRequest,
		Message:    http.StatusText(http.StatusBadRequest),
	}
	if err != nil {
		errors, ok := err.(validator.ValidationErrors)
		if !ok {
			ans.Errors = append(ans.Errors, err.Error())
		} else {
			for i := range errors {
				ans.Errors = append(ans.Errors, errors[i].Error())
			}
		}
	}
	return ans
}
