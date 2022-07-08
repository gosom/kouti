package web

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

type BaseHandler struct {
	Logger zerolog.Logger
}

func (o BaseHandler) Log(r *http.Request, format string, v ...any) {
	var ev *zerolog.Event
	ev = o.Logger.Debug()
	ev = ev.
		Str("request-id", o.GetReqID(r)).
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("query", r.URL.RawQuery).
		Str("ip", r.RemoteAddr).
		Str("user-agent", r.UserAgent())
	ev.Msgf(format, v...)
}

func (o BaseHandler) GetReqID(r *http.Request) string {
	return middleware.GetReqID(r.Context())
}

func (o BaseHandler) URLParamString(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

func (o BaseHandler) UrlParamInt(r *http.Request, key string) (int, error) {
	v := chi.URLParam(r, key)
	ans, err := strconv.Atoi(v)
	if err != nil {
		return ans, ErrHTTP{http.StatusBadRequest, err.Error()}
	}
	return ans, nil
}

func (o BaseHandler) BindJSON(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func (o BaseHandler) Json(w http.ResponseWriter, statusCode int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if value == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(value); err != nil {
		panic(err)
	}
}
