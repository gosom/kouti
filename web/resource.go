package web

import (
	"context"
	"net/http"
)

type IResourceService[M, PostP, R any] interface {
	CreateResource(ctx context.Context, p PostP) (M, error)
	GetResourceByID(ctx context.Context, id int) (M, error)
	DeleteResourceByID(ctx context.Context, id int) error
	SelectResources(ctx context.Context) ([]M, error)
	ModelToResource(ctx context.Context, m M) (R, error)
}

type ResourceHandler[M, PostP, R any] struct {
	BaseHandler
	Srv IResourceService[M, PostP, R]
}

// Get GET
func (h ResourceHandler[M, PostP, R]) Get(w http.ResponseWriter, r *http.Request) {
	resourceId, err := h.UrlParamInt(r, "id")
	if err != nil {
		ae := NewBadRequestError("cannot fetch id from url")
		h.Json(w, ae.StatusCode, ae)
		return
	}
	m, err := h.Srv.GetResourceByID(r.Context(), resourceId)
	if err != nil {
		ae := NewInternalServerError(err)
		h.Json(w, ae.StatusCode, ae)
		return
	}
	resp, err := h.Srv.ModelToResource(r.Context(), m)
	if err != nil {
		ae := NewInternalServerError(err)
		h.Json(w, ae.StatusCode, ae)
		return
	}
	h.Json(w, http.StatusOK, resp)
}

// Post
func (h ResourceHandler[M, PostP, R]) Post(w http.ResponseWriter, r *http.Request) {
	var payload PostP
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
	m, err := h.Srv.CreateResource(r.Context(), payload)
	if err != nil {
		ae := NewInternalServerError(err)
		h.Json(w, ae.StatusCode, ae)
		return
	}
	resp, err := h.Srv.ModelToResource(r.Context(), m)
	if err != nil {
		ae := NewInternalServerError(err)
		h.Json(w, ae.StatusCode, ae)
		return
	}
	h.Json(w, http.StatusCreated, resp)
}

// Delete
func (h ResourceHandler[M, PostP, R]) Delete(w http.ResponseWriter, r *http.Request) {
	resourceId, err := h.UrlParamInt(r, "id")
	if err != nil {
		ae := NewBadRequestError("cannot fetch id from url")
		h.Json(w, ae.StatusCode, ae)
		return
	}
	if err := h.Srv.DeleteResourceByID(r.Context(), resourceId); err != nil {
		ae := NewInternalServerError(err)
		h.Json(w, ae.StatusCode, ae)
		return
	}
	h.Json(w, http.StatusOK, nil)
}

// Put
func (h ResourceHandler[M, PostP, R]) Put(w http.ResponseWriter, r *http.Request) {
}

// Patch
func (h ResourceHandler[M, PostP, R]) Patch(w http.ResponseWriter, r *http.Request) {
}

// Select
func (h ResourceHandler[M, PostP, R]) Select(w http.ResponseWriter, r *http.Request) {
	items, err := h.Srv.SelectResources(r.Context())
	if err != nil {
		ae := NewInternalServerError(err)
		h.Json(w, ae.StatusCode, ae)
		return
	}
	resps := make([]R, len(items))
	for i := range resps {
		resps[i], err = h.Srv.ModelToResource(r.Context(), items[i])
		if err != nil {
			ae := NewInternalServerError(err)
			h.Json(w, ae.StatusCode, ae)
			return
		}
	}
	h.Json(w, http.StatusOK, resps)
}

// Search
func (h ResourceHandler[M, PostP, R]) Search(w http.ResponseWriter, r *http.Request) {
}
