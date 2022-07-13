package rest

import (
	"context"
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
	PageSize   int    `schema:"pageSize" validate:"omitempty,min=1,max=100"`
	SearchTerm string `schema:"searchTerm" validate:"omitempty,min=4,max=100"`
}

type User struct {
	ID        int       `json:"id"`
	Fname     string    `json:"firstName"`
	Lname     string    `json:"lastName"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type UserCreate struct {
	Fname    string `json:"firstName" validate:"required,min=4,max=100"`
	Lname    string `json:"lastName" validate:"required,min=4,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

func NewUserHandler(log zerolog.Logger, srv *UserSrv) web.ResourceHandler[UserCreate, QueryParams, User] {
	h := web.ResourceHandler[UserCreate, QueryParams, User]{}
	h.Logger = log
	h.Srv = srv
	return h
}

type UserSrv struct {
	Log zerolog.Logger
	DB  db.DB
}

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
		_, err := q.CreateUser(ctx, orm.CreateUserParams{
			Fname:  gofakeit.FirstName(),
			Lname:  gofakeit.LastName(),
			Email:  email,
			Passwd: gofakeit.Password(true, true, true, true, false, 10),
		})
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

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
