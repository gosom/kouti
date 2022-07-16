package web

import (
	"context"
	"net/http"
)

type IAuthService[Q, R any] interface {
	PerformLogin(ctx context.Context, p Q) (R, error)
}

type AuthHandler[Q, R any] struct {
	BaseHandler
	Srv IAuthService[Q, R]
}

// Login
func (h AuthHandler[Q, R]) Login(w http.ResponseWriter, r *http.Request) {
	var payload Q
	if err := h.BindJSON(r, &payload); err != nil {
		ae := NewBadRequestError("")
		h.Json(w, ae.StatusCode, ae)
		return
	}
	if err := h.Validate(payload); err != nil {
		ae := NewValidationError(err)
		h.Json(w, ae.StatusCode, ae)
		return
	}
	resp, err := h.Srv.PerformLogin(r.Context(), payload)
	if err != nil {
		ae := NewErrHTTPFromError(err)
		h.Json(w, ae.StatusCode, ae)
		return
	}
	h.Json(w, http.StatusOK, resp)
}
