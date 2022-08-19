package web

import (
	"context"
	"net/http"
)

type IResourceService[Q, P, R any] interface {
	CreateResource(ctx context.Context, p Q) (R, error)
	GetResourceByID(ctx context.Context, id string) (R, error)
	DeleteResourceByID(ctx context.Context, id string) error
	SelectResources(ctx context.Context, p P) ([]R, error)
	SearchResources(ctx context.Context, p P) ([]R, error)
}

type ResourceHandler[Q, P, R any] struct {
	BaseHandler
	Srv IResourceService[Q, P, R]
}

// Get GET
func (h ResourceHandler[Q, P, R]) Get(w http.ResponseWriter, r *http.Request) {
	resourceId := h.URLParamString(r, "id")
	if len(resourceId) == 0 {
		ae := NewBadRequestError("cannot fetch id from url")
		h.Json(w, ae.StatusCode, ae)
		return
	}
	resp, err := h.Srv.GetResourceByID(r.Context(), resourceId)
	if err != nil {
		ae := NewErrHTTPFromError(err)
		h.Json(w, ae.StatusCode, ae)
		return
	}
	h.Json(w, http.StatusOK, resp)
}

// Post
func (h ResourceHandler[Q, P, R]) Post(w http.ResponseWriter, r *http.Request) {
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
	resp, err := h.Srv.CreateResource(r.Context(), payload)
	if err != nil {
		h.Logger.Error().AnErr("error", err).Msg("on CreateResource")
		ae := NewErrHTTPFromError(err)
		h.Json(w, ae.StatusCode, ae)
		return
	}
	h.Json(w, http.StatusCreated, resp)
}

// Delete
func (h ResourceHandler[Q, P, R]) Delete(w http.ResponseWriter, r *http.Request) {
	resourceId := h.URLParamString(r, "id")
	if len(resourceId) == 0 {
		ae := NewBadRequestError("cannot fetch id from url")
		h.Json(w, ae.StatusCode, ae)
		return
	}
	if err := h.Srv.DeleteResourceByID(r.Context(), resourceId); err != nil {
		ae := NewInternalServerError(err)
		h.Json(w, ae.StatusCode, ae)
		return
	}
	h.Json(w, http.StatusNoContent, nil)
}

// Put
func (h ResourceHandler[Q, P, R]) Put(w http.ResponseWriter, r *http.Request) {
}

// Patch
func (h ResourceHandler[Q, P, R]) Patch(w http.ResponseWriter, r *http.Request) {
}

// Select
func (h ResourceHandler[Q, P, R]) Select(w http.ResponseWriter, r *http.Request) {
	var qp P
	if err := h.BindQueryParams(r, &qp); err != nil {
		ae := NewBadRequestError("")
		h.Json(w, ae.StatusCode, ae)
		return
	}
	if err := h.Validate(qp); err != nil {
		ae := NewValidationError(err)
		h.Json(w, ae.StatusCode, ae)
		return
	}
	items, err := h.Srv.SelectResources(r.Context(), qp)
	if err != nil {
		ae := NewInternalServerError(err)
		h.Json(w, ae.StatusCode, ae)
		return
	}
	h.Json(w, http.StatusOK, items)
}

// Search
func (h ResourceHandler[Q, P, R]) Search(w http.ResponseWriter, r *http.Request) {
	var qp P
	if err := h.BindQueryParams(r, &qp); err != nil {
		ae := NewBadRequestError("")
		h.Json(w, ae.StatusCode, ae)
		return
	}
	if err := h.Validate(qp); err != nil {
		ae := NewValidationError(err)
		h.Json(w, ae.StatusCode, ae)
		return
	}
	items, err := h.Srv.SearchResources(r.Context(), qp)
	if err != nil {
		ae := NewInternalServerError(err)
		h.Json(w, ae.StatusCode, ae)
		return
	}
	h.Json(w, http.StatusOK, items)
}
