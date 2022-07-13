package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
)

type ErrHTTP struct {
	StatusCode int      `json:"-"`
	Message    string   `json:"message,omitempty"`
	Errors     []string `json:"errors,omitempty"`
}

func (e ErrHTTP) Error() string {
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

func NewNotFoundError() error {
	ans := ErrHTTP{
		StatusCode: http.StatusNotFound,
		Message:    "resource not found",
	}
	return &ans
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

func NewInternalServerError(err error) ErrHTTP {
	ans := ErrHTTP{
		StatusCode: http.StatusInternalServerError,
		Message:    err.Error(),
	}
	return ans
}

func NewErrHTTPFromError(err error) ErrHTTP {
	switch v := err.(type) {
	case *ErrHTTP:
		return *v
	case ErrHTTP:
		return v
	}
	e := ErrHTTP{
		StatusCode: http.StatusInternalServerError,
		Message:    http.StatusText(http.StatusInternalServerError),
	}
	var pgErr *pgconn.PgError
	if errors.Is(err, pgx.ErrNoRows) {
		e.StatusCode = http.StatusNotFound
		e.Message = "resource not found"
	} else if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			e.StatusCode = http.StatusConflict
			e.Message = "resource already exists"
		}
	}
	return e
}
