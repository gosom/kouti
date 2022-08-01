package um

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/gosom/kouti/dbdriver"
	"github.com/rs/zerolog"
)

type Config struct {
	Log         zerolog.Logger
	DB          dbdriver.DB
	SystemRoles []string
}

type Service struct {
	logger      zerolog.Logger
	db          dbdriver.DB
	systemRoles []string
	Schema      ISchema
	Roles       IRoleRepo
	Users       IUserRepo
}

func NewService(cfg Config) *Service {
	ans := Service{
		logger:      cfg.Log,
		db:          cfg.DB,
		systemRoles: cfg.SystemRoles,
		Schema: &Schema{
			logger: cfg.Log,
		},
		Roles: &RoleRepo{
			logger: cfg.Log,
		},
		Users: &UserRepo{
			logger: cfg.Log,
		},
	}
	return &ans
}

func (s *Service) InitSchema(ctx context.Context) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := s.Schema.CreateTables(ctx, tx); err != nil {
		return err
	}
	if len(s.systemRoles) > 0 {
		if _, err := s.Roles.Insert(ctx, tx, s.systemRoles...); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

type RegisterUserOpts struct {
	Identity string
	Passwd   string
	Roles    []string

	BeforeCommit func(ctx context.Context, conn dbdriver.DBTX, u User) error
	AfterCommit  func(ctx context.Context, conn dbdriver.DBTX, u User) error
}

func (s *Service) RegisterUser(ctx context.Context, p RegisterUserOpts) (User, error, error) {
	if len(p.Roles) == 0 {
		return User{}, errors.New("you must provide at least one role for the user"), nil
	}
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return User{}, err, nil
	}
	defer tx.Rollback(ctx)

	roles, err := s.Roles.Select(ctx, tx, RoleSelectParams{
		Names:    p.Roles,
		UseNames: true,
	})
	if err != nil {
		return User{}, err, nil
	}
	if len(p.Roles) != len(roles) {
		return User{}, errors.New("roles do not exist"), nil
	}

	u, err := s.Users.Insert(ctx, tx, InsertUserParams{
		Identity: p.Identity,
		Passwd:   p.Passwd,
	})
	if err != nil {
		return User{}, err, nil
	}

	if err := s.Roles.AssignRolesToUser(ctx, tx, UserRoleAssignParams{
		UserID: u.ID,
		Roles:  roles,
	}); err != nil {
		return User{}, err, nil
	}
	u.Roles = roles
	if p.BeforeCommit != nil {
		if err := p.BeforeCommit(ctx, tx, u); err != nil {
			return u, err, nil
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return u, err, nil
	}

	if p.AfterCommit != nil {
		if err := p.AfterCommit(ctx, tx, u); err != nil {
			return u, nil, err
		}
	}
	return u, nil, nil
}

func (s *Service) GetUserByUID(ctx context.Context, uid string) (User, error) {
	_uid, err := uuid.Parse(uid)
	if err != nil {
		return User{}, err
	}
	p := GetUserParams{
		UseUID: true,
		UID:    _uid,
	}
	return s.Users.Get(ctx, s.db.RawConn(), p)
}

func (s *Service) SelectUsers(ctx context.Context, p SelectUserParams) ([]User, error) {
	return s.Users.Select(ctx, s.db.RawConn(), p)
}

func (s *Service) DeleteUser(ctx context.Context, p GetUserParams) error {
	return s.Users.Delete(ctx, s.db.RawConn(), p)
}

func (s *Service) Login(ctx context.Context, u string, p string) (User, error) {
	param := UserLoginParams{
		Identity: u,
		Passwd:   p,
	}
	return s.Users.Login(ctx, s.db.RawConn(), param)
}
