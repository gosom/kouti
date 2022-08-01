package um

import (
	"context"

	"github.com/gosom/kouti/dbdriver"
	"github.com/rs/zerolog"
)

type ISchema interface {
	CreateTables(ctx context.Context, conn dbdriver.DBTX) error
}

type Schema struct {
	logger zerolog.Logger
}

func (s *Schema) CreateTables(ctx context.Context, conn dbdriver.DBTX) error {
	queries := []string{
		createUsrTblQ,
		createRoleTblQ,
		createUserRolesTblQ,
		setTimestampFnQ,
		usersDropTsTrigger,
		usersSetTsTrigger,
	}
	for _, q := range queries {
		if _, err := conn.Exec(ctx, q); err != nil {
			return err
		}
	}
	return nil
}
