package um

import (
	"context"

	"github.com/gosom/kouti/dbdriver"
	"github.com/rs/zerolog"
)

type IRoleRepo interface {
	Insert(ctx context.Context, conn dbdriver.DBTX, name ...string) ([]Role, error)
	Get(ctx context.Context, conn dbdriver.DBTX, id int) (Role, error)
	Select(ctx context.Context, conn dbdriver.DBTX, p RoleSelectParams) ([]Role, error)
	Delete(ctx context.Context, conn dbdriver.DBTX, id int) error
	Update(ctx context.Context, conn dbdriver.DBTX, p RoleUpdateParams) (Role, error)
	AssignRolesToUser(ctx context.Context, conn dbdriver.DBTX, p UserRoleAssignParams) error
}

type UserRoleAssignParams struct {
	UserID int
	Roles  []Role
}

type RoleSelectParams struct {
	Ids    []int
	UseIds bool

	Names    []string
	UseNames bool
}

type RoleUpdateParams struct {
	ID   int
	Name string
}

type RoleRepo struct {
	logger zerolog.Logger
}

func (r *RoleRepo) Insert(ctx context.Context, conn dbdriver.DBTX, name ...string) ([]Role, error) {
	rows, err := conn.Query(ctx, insertRolesQ, name)
	if err != nil {
		return nil, err
	}
	items := make([]Role, 0, len(name))
	for rows.Next() {
		var r Role
		if err := rows.Scan(&r.ID, &r.Name, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, err
}

func (r *RoleRepo) Get(ctx context.Context, conn dbdriver.DBTX, id int) (Role, error) {
	// TODO
	return Role{}, nil
}

func (r *RoleRepo) Select(ctx context.Context, conn dbdriver.DBTX, p RoleSelectParams) ([]Role, error) {
	rows, err := conn.Query(ctx, selectRolesQ, p.UseIds, p.Ids, p.UseNames, p.Names)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Role
	for rows.Next() {
		var role Role
		if err := rows.Scan(&role.ID, &role.Name, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, role)
	}
	if err := rows.Err(); err != nil {
		return items, err
	}
	return items, nil
}

func (r *RoleRepo) Delete(ctx context.Context, conn dbdriver.DBTX, id int) error {
	// TODO
	return nil
}

func (r *RoleRepo) Update(ctx context.Context, conn dbdriver.DBTX, p RoleUpdateParams) (Role, error) {
	// TODO
	return Role{}, nil
}
func (r *RoleRepo) AssignRolesToUser(ctx context.Context, conn dbdriver.DBTX, p UserRoleAssignParams) error {
	roleIds := make([]int, len(p.Roles), len(p.Roles))
	for i := range p.Roles {
		roleIds[i] = p.Roles[i].ID
	}
	_, err := conn.Exec(ctx, insertUserRolesQ, p.UserID, roleIds)
	return err
}
