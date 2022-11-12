package web

import "net/http"

type IBaseHandler interface {
	Validate(s any, fields ...string) error
	Log(r *http.Request, format string, v ...any)
	GetReqID(r *http.Request) string
	URLParamString(r *http.Request, key string) string
	URLParamInt(r *http.Request, key string) (int, error)
	BindJSON(r *http.Request, v any) error
	BindJSONValidate(r *http.Request, v any) error
	BindQueryParams(r *http.Request, v any) error
	Json(w http.ResponseWriter, statusCode int, value any)
}
