package todos

import (
	"context"

	"github.com/rs/zerolog"
)

type TodoSrv struct {
	Log zerolog.Logger
}

// CreateResource Creates a todo resource
// @summary Create a todo
// @id create-todo
// @produce json
// @param Body body TodoCreate true "the body to create a todo"
// @success 201 {object} Todo
// @failure 400 {object} web.ErrHTTP
// @failure 409 {object} web.ErrHTTP
// @failure 500 {object} web.ErrHTTP
// @router /users [post]
func (o *TodoSrv) CreateResource(ctx context.Context, p TodoCreate) (Todo, error) {
	var ans Todo
	// TODO
	return ans, nil
}

// GetResourceByID Returns a specific user
// @summary Returns user with id
// @id get-user
// @produce json
// @param id path string true "the id of the user to fetch"
// @success 200 {object} Todo
// @failure 400 {object} web.ErrHTTP
// @failure 500 {object} web.ErrHTTP
// @router /users/{id} [get]
func (o *TodoSrv) GetResourceByID(ctx context.Context, id string) (Todo, error) {
	var ans Todo
	// TODO
	return ans, nil
}

// DeleteResourceByID Deletes a specific user
// @summary Deletes user with id
// @id delete-user
// @produce json
// @param id path string true "the id of the user to delete"
// @success 204
// @failure 400 {object} web.ErrHTTP
// @failure 500 {object} web.ErrHTTP
// @router /users/{id} [delete]
func (o *TodoSrv) DeleteResourceByID(ctx context.Context, id string) error {
	// TODO
	return nil
}

// SelectResources select an array of users meeting criteria
// @summary Returns a list of users
// @id select-users
// @produce json
// @param next query int false "the id of the next user (used for pagination)"
// @param pageSize query int false "the number of results per page"
// @param email query string false "filter by email"
// @param firstName query string false "filter by firstName"
// @param lastName query string false "filter by lastName"
// @success 200 {array} Todo
// @failure 400 {object} web.ErrHTTP
// @failure 500 {object} web.ErrHTTP
// @router /users/ [get]
func (o *TodoSrv) SelectResources(ctx context.Context, qp TodoQueryParams) ([]Todo, error) {
	// TODO
	return nil, nil
}

// SearchResources searches for users that match the searchTerm
// @summary Returns a list of users
// @id search-users
// @produce json
// @param next query int false "the id of the next user (used for pagination)"
// @param pageSize query int false "the number of results per page"
// @param searchTerm query string false "search term"
// @success 200 {array} Todo
// @failure 400 {object} web.ErrHTTP
// @failure 500 {object} web.ErrHTTP
// @router /users/search [get]
func (o *TodoSrv) SearchResources(ctx context.Context, qp TodoQueryParams) ([]Todo, error) {
	// TODO
	return nil, nil
}
