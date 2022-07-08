package web

import "net/http"

type DefaultHandler struct {
	BaseHandler
}

func (d DefaultHandler) NotFound(w http.ResponseWriter, r *http.Request) {
	e := ErrHTTP{http.StatusNotFound, http.StatusText(http.StatusNotFound)}
	d.Json(w, e.StatusCode, e)
}

func (d DefaultHandler) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	e := ErrHTTP{http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed)}
	d.Json(w, e.StatusCode, e)
}
