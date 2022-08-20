package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/gosom/kouti/dbdriver"
	"github.com/rs/zerolog"
)

type SelectParams struct {
	UseContent     bool
	Content        string
	UseIsCompleted bool
	IsCompleted    bool
}

type Store interface {
	Create(ctx context.Context, todo Todo) error
	Get(ctx context.Context, id uuid.UUID) (Todo, error)
	Update(ctx context.Context, todo Todo) error
	Delete(ctx context.Context, id uuid.UUID) error
	Select(ctx context.Context, p SelectParams) ([]Todo, error)
}

type postgresStore struct {
	log zerolog.Logger
	db  *dbdriver.DB
}

func NewPostgresStore(log zerolog.Logger, conn *dbdriver.DB) *postgresStore {
	ans := postgresStore{
		log: log,
		db:  conn,
	}
	return &ans
}

func (s *postgresStore) Create(ctx context.Context, todo Todo) error {
	const q = `insert into todos
	(id, content, is_completed, created_at, updated_at)
	values ($1, $2, $3, $4, $5)`
	_, err := s.db.RawConn().Exec(
		ctx,
		q,
		todo.ID, todo.Content, todo.IsCompleted, todo.CreatedAt, todo.UpdatedAt,
	)

	return err
}
func (s *postgresStore) Get(ctx context.Context, id uuid.UUID) (Todo, error) {
	const q = `select * from todos where id = $1`
	row := s.db.RawConn().QueryRow(ctx, q, id)
	return rowToTodo(row)
}

func (s *postgresStore) Update(ctx context.Context, todo Todo) error {
	const q = `update todos set
	content = $1, is_completed = $2, updated_at = $3
	where id = $4
	`
	_, err := s.db.RawConn().Exec(
		ctx,
		q,
		todo.Content, todo.IsCompleted, todo.UpdatedAt, todo.ID,
	)
	return err
}

func (s *postgresStore) Delete(ctx context.Context, id uuid.UUID) error {
	const q = `delete from todos where id = $1`
	_, err := s.db.RawConn().Exec(
		ctx,
		q,
		id,
	)
	return err
}

func (s *postgresStore) Select(ctx context.Context, p SelectParams) ([]Todo, error) {
	const q = `select * from todos
	where
	(CASE WHEN $1::boolean THEN content ilike 
		concat('%', $2::text, '%') ELSE true END) AND
	(CASE WHEN $3::boolean THEN is_completed = $4 ELSE true END)
	`
	rows, err := s.db.RawConn().Query(
		ctx,
		q,
		p.UseContent, p.Content,
		p.UseIsCompleted, p.IsCompleted,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Todo
	for rows.Next() {
		todo, err := rowToTodo(rows)
		if err != nil {
			return items, err
		}
		items = append(items, todo)
	}
	if len(items) == 0 {
		items = []Todo{}
	}
	return items, nil
}

func rowToTodo(r dbdriver.RowScanner) (ans Todo, err error) {
	err = r.Scan(
		&ans.ID, &ans.Content, &ans.IsCompleted, &ans.CreatedAt, &ans.UpdatedAt,
	)
	return
}
