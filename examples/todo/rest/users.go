package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"

	"github.com/gosom/kouti/examples/todo/db"
	"github.com/gosom/kouti/examples/todo/orm"
	"github.com/gosom/kouti/web"

	"github.com/brianvoe/gofakeit/v6"
)

type QueryParams struct {
	Next       int    `schema:"next"`
	Email      string `schema:"email" validate:"omitempty,email"`
	Fname      string `schema:"firstName" validate:"omitempty,min=1,max=100"`
	Lname      string `schema:"lastName" validate:"omitempty,min=1,max=100"`
	PageSize   int    `schema:"pageSize" validate:"omitempty,min=1,max=100"`
	SearchTerm string `schema:"searchTerm" validate:"omitempty,min=4,max=100"`
}

type User struct {
	ID        int       `json:"id" example:"1"`
	Fname     string    `json:"firstName" example:"Aris"`
	Lname     string    `json:"lastName" example:"Paparis"`
	Email     string    `json:"email" example:"aris.paparis@example.com"`
	CreatedAt time.Time `json:"created_at" example:"2022-07-16T00:53:16.535668Z" format:"date-time"`
}

type UserCreate struct {
	Fname    string `json:"firstName" validate:"required,min=4,max=100" example:"Aris"`
	Lname    string `json:"lastName" validate:"required,min=4,max=100" example:"Paparis"`
	Email    string `json:"email" validate:"required,email" example:"aris.paparis@example.com"`
	Password string `json:"password" validate:"required,password" example:"Ar9Sp7891!!#"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required,email" example:"aris.paparis@example.com"`
	Password string `json:"password" validate:"required" example:"Ar9Sp7891!!#"`
}

type L struct {
	AccessToken string `json:"accessToken"`
}

func NewAuthHandler(log zerolog.Logger, srv *UserSrv) web.AuthHandler[UserLogin, L] {
	h := web.AuthHandler[UserLogin, L]{}
	h.Logger = log
	h.Srv = srv
	return h
}

func NewUserHandler(log zerolog.Logger, srv *UserSrv) web.ResourceHandler[UserCreate, QueryParams, User] {
	h := web.ResourceHandler[UserCreate, QueryParams, User]{}
	h.Logger = log
	h.Srv = srv
	return h
}

type IAuth interface {
	GetAccessToken(u any) (string, error)
}

type UserSrv struct {
	Log  zerolog.Logger
	DB   db.DB
	Auth IAuth
}

// TODO if swag was supporting generics this code could have been moved

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
	m, err := o.DB.CreateUser(ctx, orm.CreateUserParams{
		Fname:  p.Fname,
		Lname:  p.Lname,
		Email:  p.Email,
		Passwd: p.Password,
	})
	if err != nil {
		return ans, err
	}
	ans.ID = int(m.ID)
	ans.Fname = m.Fname
	ans.Lname = m.Lname
	ans.Email = m.Email
	ans.CreatedAt = m.CreatedAt
	return ans, nil
}

// GetResourceByID Returns a specific user
// @summary Returns user with id
// @id get-user
// @produce json
// @param id path int true "the id of the user to fetch"
// @success 200 {object} User
// @failure 400 {object} web.ErrHTTP
// @failure 500 {object} web.ErrHTTP
// @router /users/{id} [get]
func (o *UserSrv) GetResourceByID(ctx context.Context, id int) (User, error) {
	var ans User
	m, err := o.DB.GetUserByID(ctx, int32(id))
	if err != nil {
		return ans, err
	}
	ans.ID = int(m.ID)
	ans.Fname = m.Fname
	ans.Lname = m.Lname
	ans.Email = m.Email
	ans.CreatedAt = m.CreatedAt
	return ans, nil
}

// DeleteResourceByID Deletes a specific user
// @summary Deletes user with id
// @id delete-user
// @produce json
// @param id path int true "the id of the user to fetch"
// @success 204
// @failure 400 {object} web.ErrHTTP
// @failure 500 {object} web.ErrHTTP
// @router /users/{id} [delete]
func (o *UserSrv) DeleteResourceByID(ctx context.Context, id int) error {
	n, err := o.DB.DeleteUserByID(ctx, int32(id))
	if err != nil {
		return err
	}
	if n == 0 {
		return web.NewNotFoundError()
	}
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
func (o *UserSrv) SelectResources(ctx context.Context, qp QueryParams) ([]User, error) {
	p := orm.ListUsersParams{
		UseRlimit: true,
	}
	if qp.PageSize > 0 {
		p.Rlimit = int32(qp.PageSize)
	}
	if qp.Next > 0 {
		p.ID = int32(qp.Next)
	}
	if qp.Email != "" {
		p.WhereEmail = true
		p.Email = qp.Email
	}
	if qp.Fname != "" {
		p.WhereFname = true
		p.Fname = qp.Fname
	}
	if qp.Fname != "" {
		p.WhereLname = true
		p.Lname = qp.Lname
	}
	rows, err := o.DB.ListUsers(ctx, p)
	if err != nil {
		return nil, err
	}
	var items []User
	if len(rows) == 0 {
		items = []User{}
	} else {
		items = make([]User, len(rows))
		for i := range rows {
			items[i] = User{
				ID:        int(rows[i].ID),
				Fname:     rows[i].Fname,
				Lname:     rows[i].Lname,
				Email:     rows[i].Email,
				CreatedAt: rows[i].CreatedAt,
			}
		}
	}
	return items, nil
}

func (o *UserSrv) InsertRandomUsers(ctx context.Context, num int) error {
	tx, err := o.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	q := o.DB.WithTx(tx)
	seen := make(map[string]bool)

	for i := 0; i < num; i++ {
		var email string
		for {
			email = gofakeit.Email()
			if !seen[email] {
				seen[email] = true
				break
			}
		}
		u, err := q.CreateUser(ctx, orm.CreateUserParams{
			Fname:  gofakeit.FirstName(),
			Lname:  gofakeit.LastName(),
			Email:  email,
			Passwd: gofakeit.Password(true, true, true, true, false, 10),
		})
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx,
			`INSERT INTO casbin(id, ptype, v0, v1, v2, v3, v4, v5) VALUES($1, $2, $3, $4, $5, $6, $7, $8)`,
			i+1, "p", fmt.Sprintf("%d", u.ID),
			fmt.Sprintf("/api/v1/users/%d", u.ID), "(GET)|(DELETE)", "", "", "",
		)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
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
func (o *UserSrv) SearchResources(ctx context.Context, qp QueryParams) ([]User, error) {
	p := db.UserFtsParams{
		Phrase: qp.SearchTerm,
	}
	if qp.PageSize > 0 {
		p.Limit = qp.PageSize
	}
	if qp.Next > 0 {
		p.ID = qp.Next
	}
	rows, err := o.DB.SearchUsers(ctx, p)
	if err != nil {
		return nil, err
	}
	var items []User
	if len(rows) == 0 {
		items = []User{}
	} else {
		items = make([]User, len(rows))
		for i := range rows {
			items[i] = User{
				ID:        int(rows[i].ID),
				Fname:     rows[i].Fname,
				Lname:     rows[i].Lname,
				Email:     rows[i].Email,
				CreatedAt: rows[i].CreatedAt,
			}
		}
	}
	return items, nil
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
	qp := db.UserLoginParams{
		Email:  p.Email,
		Passwd: p.Password,
	}
	u, err := o.DB.GetUserByEmailPasswd(ctx, qp)
	if err != nil {
		ae := web.ErrHTTP{
			StatusCode: http.StatusUnauthorized,
			Message:    http.StatusText(http.StatusUnauthorized),
		}
		return L{}, ae
	}
	accessToken, err := o.Auth.GetAccessToken(u.ID)
	if err != nil {
		ae := web.ErrHTTP{
			StatusCode: http.StatusUnauthorized,
			Message:    http.StatusText(http.StatusUnauthorized),
		}
		return L{}, ae
	}
	return L{AccessToken: accessToken}, nil
}
