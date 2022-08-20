package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gosom/kouti/web"
	"github.com/rs/zerolog"
)

type TodoCreate struct {
	Content     string `json:"content" validate:"required,min=4,max=140"`
	IsCompleted bool   `json:"is_completed"`
}

type TodoSelectParams struct {
	Content     *string `schema:"content" validate:"omitempty,min=1,max=140"`
	IsCompleted *bool   `schema:"is_completed"`
}

// TodoHandler
type TodoHandler struct {
	web.BaseHandler
	store Store
}

func NewToDoHandler(log zerolog.Logger, s Store) TodoHandler {
	ans := TodoHandler{}
	ans.Logger = log
	ans.store = s
	return ans
}

// Post Creates a todo
// @summary Create a todo
// @id create-todo
// @produce json
// @param Body body TodoCreate true "the body to create a todo"
// @success 201 {object} Todo
// @failure 400 {object} web.ErrHTTP
// @failure 500 {object} web.ErrHTTP
// @router /todos [post]
func (h TodoHandler) Post(w http.ResponseWriter, r *http.Request) {
	var payload TodoCreate
	if err := h.BindJSONValidate(r, &payload); err != nil {
		ae := web.NewBadRequestError("")
		h.Json(w, ae.StatusCode, ae)
		return
	}

	now := time.Now().UTC()
	todo := Todo{
		ID:          uuid.New(),
		Content:     payload.Content,
		IsCompleted: payload.IsCompleted,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := h.store.Create(r.Context(), todo); err != nil {
		ae := web.NewInternalServerError(err)
		h.Json(w, ae.StatusCode, ae)
		return
	}
	h.Json(w, http.StatusCreated, todo)
}

func (h TodoHandler) getUuidFromReq(r *http.Request) (uuid.UUID, error) {
	id := h.URLParamString(r, "id")
	uid, err := uuid.Parse(id)
	return uid, err
}

// Get Returns a specific todo
// @summary Returns todo with id
// @id get-todo
// @produce json
// @param id path string true "the id of the todo to fetch"
// @success 200 {object} Todo
// @failure 400 {object} web.ErrHTTP
// @failure 500 {object} web.ErrHTTP
// @router /todos/{id} [get]
func (h TodoHandler) Get(w http.ResponseWriter, r *http.Request) {
	uid, err := h.getUuidFromReq(r)
	if err != nil {
		ae := web.NewBadRequestError(err.Error())
		h.Json(w, ae.StatusCode, ae)
		return
	}
	todo, err := h.store.Get(r.Context(), uid)
	if err != nil { // this may be a not found error TODO
		ae := web.NewInternalServerError(err)
		h.Json(w, ae.StatusCode, ae)
		return
	}
	h.Json(w, http.StatusOK, todo)
}

// Put Updates a specific todo
// @summary Updates todo with id
// @id update-todo
// @produce json
// @param id path string true "the id of the todo to update"
// @success 204
// @failure 400 {object} web.ErrHTTP
// @failure 500 {object} web.ErrHTTP
// @router /todos/{id} [put]
func (h TodoHandler) Put(w http.ResponseWriter, r *http.Request) {
	uid, err := h.getUuidFromReq(r)
	if err != nil {
		ae := web.NewBadRequestError(err.Error())
		h.Json(w, ae.StatusCode, ae)
		return
	}
	var payload TodoCreate
	if err := h.BindJSONValidate(r, &payload); err != nil {
		ae := web.NewBadRequestError(err.Error())
		h.Json(w, ae.StatusCode, ae)
		return
	}

	todo := Todo{
		ID:          uid,
		Content:     payload.Content,
		IsCompleted: payload.IsCompleted,
		UpdatedAt:   time.Now().UTC(),
	}

	if err := h.store.Update(r.Context(), todo); err != nil {
		ae := web.NewInternalServerError(err)
		h.Json(w, ae.StatusCode, ae)
		return
	}
	h.Json(w, http.StatusNoContent, nil)
}

// Delete deletes a specific todo
// @summary deletes todo with id
// @id delete-todo
// @produce json
// @param id path string true "the id of the todo to delete"
// @success 204
// @failure 400 {object} web.ErrHTTP
// @failure 500 {object} web.ErrHTTP
// @router /todos/{id} [delete]
func (h TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	uid, err := h.getUuidFromReq(r)
	if err != nil {
		ae := web.NewBadRequestError(err.Error())
		h.Json(w, ae.StatusCode, ae)
		return
	}
	err = h.store.Delete(r.Context(), uid)
	if err != nil { // this may be a not found error TODO
		ae := web.NewInternalServerError(err)
		h.Json(w, ae.StatusCode, ae)
		return
	}
	h.Json(w, http.StatusNoContent, nil)
}

// Select select an array of todos meeting criteria
// @summary Returns a list of todos
// @id select-todos
// @produce json
// @param content query string false "returns todos containing Content"
// @param is_completed query string false "t to return is_completed todos and f to with completed false"
// @success 200 {array} Todo
// @failure 400 {object} web.ErrHTTP
// @failure 500 {object} web.ErrHTTP
// @router /todos/ [get]
func (h TodoHandler) Select(w http.ResponseWriter, r *http.Request) {
	var params TodoSelectParams
	if err := h.BindQueryParams(r, &params); err != nil {
		ae := web.NewBadRequestError(err.Error())
		h.Json(w, ae.StatusCode, ae)
		return
	}
	p := SelectParams{}
	if params.Content != nil {
		p.UseContent = true
		p.Content = *params.Content
	}
	if params.IsCompleted != nil {
		p.UseIsCompleted = true
		p.IsCompleted = *params.IsCompleted
	}
	items, err := h.store.Select(r.Context(), p)
	if err != nil {
		ae := web.NewInternalServerError(err)
		h.Json(w, ae.StatusCode, ae)
		return
	}
	h.Json(w, http.StatusOK, items)
}
