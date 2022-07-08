package web

import "fmt"

type ErrHTTP struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message,omitempty"`
}

func (e ErrHTTP) Error() string {
	if e.StatusCode > 0 {
		msg := fmt.Sprintf("%d: %s", e.StatusCode, e.Message)
		return msg
	}
	return e.Message
}
