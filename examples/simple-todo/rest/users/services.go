package users

import (
	"context"

	"github.com/gosom/kouti/um"
	"github.com/rs/zerolog"
)

type UserSrv struct {
	Log zerolog.Logger
	Srv *um.Service
}

// CreateResource Creates a user resource
// @summary Create a user
// @id create-user
// @produce json
// @param Body body UserCreate true "the body to create a user"
// @success 201 {object} User
// @failure 400 {object} web.ErrHTTP
// @failure 409 {object} web.ErrHTTP
// @failure 500 {object} web.ErrHTTP
// @router /users [post]
func (o *UserSrv) CreateResource(ctx context.Context, p UserCreate) (User, error) {
	var ans User
	// TODO
	return ans, nil
}

// GetResourceByID Returns a specific user
// @summary Returns user with id
// @id get-user
// @produce json
// @param id path string true "the id of the user to fetch"
// @success 200 {object} User
// @failure 400 {object} web.ErrHTTP
// @failure 500 {object} web.ErrHTTP
// @router /users/{id} [get]
func (o *UserSrv) GetResourceByID(ctx context.Context, id string) (User, error) {
	var ans User
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
func (o *UserSrv) DeleteResourceByID(ctx context.Context, id string) error {
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
// @success 200 {array} User
// @failure 400 {object} web.ErrHTTP
// @failure 500 {object} web.ErrHTTP
// @router /users/ [get]
func (o *UserSrv) SelectResources(ctx context.Context, qp UserQueryParams) ([]User, error) {
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
// @success 200 {array} User
// @failure 400 {object} web.ErrHTTP
// @failure 500 {object} web.ErrHTTP
// @router /users/search [get]
func (o *UserSrv) SearchResources(ctx context.Context, qp UserQueryParams) ([]User, error) {
	// TODO
	return nil, nil
}

// PerformLogin login with email password
// @summary returns a JWT access token
// @id login-user
// @produce json
// @param Body body UserLogin true "the body to login a user"
// @success 200 {array} L
// @failure 400 {object} web.ErrHTTP
// @failure 500 {object} web.ErrHTTP
// @router /users/login [post]
func (o *UserSrv) PerformLogin(ctx context.Context, p UserLogin) (L, error) {
	return L{AccessToken: ""}, nil
}
