package um

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/gosom/kouti/dbdriver"
)

type InsertUserParams struct {
	Identity string
	Passwd   string
}

type GetUserParams struct {
	UseUID bool
	UID    uuid.UUID

	UseID bool
	ID    int

	UseIdentity bool
	Identity    string
}

type SelectUserParams struct {
	AfterID int

	UseLimit bool
	Limit    int
}

type UpdateUserParams struct {
}

type IUserRepo interface {
	Insert(ctx context.Context, conn dbdriver.DBTX, p InsertUserParams) (User, error)
	Get(ctx context.Context, conn dbdriver.DBTX, p GetUserParams) (User, error)
	Select(ctx context.Context, conn dbdriver.DBTX, p SelectUserParams) ([]User, error)
	Delete(ctx context.Context, conn dbdriver.DBTX, p GetUserParams) error
	Update(ctx context.Context, conn dbdriver.DBTX, p UpdateUserParams) (User, error)
}

type UserRepo struct {
	logger zerolog.Logger
}

func (u *UserRepo) Insert(ctx context.Context, conn dbdriver.DBTX, p InsertUserParams) (User, error) {
	var user User
	if err := conn.QueryRow(ctx, usersCreateQ, p.Identity, p.Passwd).Scan(
		&user.ID,
		&user.UID,
		&user.Identity,
		&user.EncPasswd,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return User{}, err
	}
	return user, nil
}

func (u *UserRepo) rowToUser(row dbdriver.RowScanner) (user User, err error) {
	err = row.Scan(
		&user.ID,
		&user.UID,
		&user.Identity,
		&user.EncPasswd,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return
}

func (u *UserRepo) Get(ctx context.Context, conn dbdriver.DBTX, p GetUserParams) (User, error) {
	row := conn.QueryRow(ctx, usersGetQ, p.UseUID, p.UID, p.UseID, p.ID, p.UseIdentity, p.Identity)
	return u.rowToUser(row)
}

func (u *UserRepo) Select(ctx context.Context, conn dbdriver.DBTX, p SelectUserParams) ([]User, error) {
	rows, err := conn.Query(ctx, usersSelectQ, p.AfterID, p.UseLimit, p.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		e, err := u.rowToUser(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, e)
	}
	return items, rows.Err()
}

func (u *UserRepo) Delete(ctx context.Context, conn dbdriver.DBTX, p GetUserParams) error {
	_, err := conn.Exec(ctx, usersDeleteQ, p.UseUID, p.UID, p.UseID, p.ID, p.UseIdentity, p.Identity)
	return err
}

func (u *UserRepo) Update(ctx context.Context, conn dbdriver.DBTX, p UpdateUserParams) (User, error) {
	// TODO
	return User{}, nil
}
