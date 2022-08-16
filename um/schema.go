package um

import (
	"context"
	"fmt"

	"github.com/gosom/kouti/dbdriver"
	"github.com/rs/zerolog"
)

type ISchema interface {
	Up(ctx context.Context, conn dbdriver.DBTX) error
	Down(ctx context.Context, conn dbdriver.DBTX) error
}

type Schema struct {
	logger zerolog.Logger
}

func (s *Schema) exec(ctx context.Context, queries []string, conn dbdriver.DBTX) error {
	for _, q := range queries {
		fmt.Println(q)
		if _, err := conn.Exec(ctx, q); err != nil {
			return err
		}
	}
	return nil
}

func (s *Schema) Up(ctx context.Context, conn dbdriver.DBTX) error {
	queries, _, err := readQueries()
	if err != nil {
		return err
	}
	return s.exec(ctx, queries, conn)
}

func (s *Schema) Down(ctx context.Context, conn dbdriver.DBTX) error {
	_, queries, err := readQueries()
	if err != nil {
		return err
	}
	return s.exec(ctx, queries, conn)
}
