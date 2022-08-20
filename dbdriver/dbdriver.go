package dbdriver

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
)

type RowScanner interface {
	Scan(dest ...interface{}) error
}

type DBTX interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
}

type Config struct {
	Logger     zerolog.Logger
	ConnString string
}

type DB struct {
	logger zerolog.Logger
	dbconn *pgxpool.Pool
}

func New(ctx context.Context, cfg Config) (*DB, error) {
	ans := DB{
		logger: cfg.Logger,
	}
	var err error
	ans.dbconn, err = Connect(ctx, cfg.ConnString)
	return &ans, err
}

func (o *DB) RawConn() *pgxpool.Pool {
	return o.dbconn
}

func (o *DB) Begin(ctx context.Context) (pgx.Tx, error) {
	return o.dbconn.Begin(ctx)
}

func (o *DB) Close() {
	o.dbconn.Close()
}

// Connect connString like:
// //   # Example DSN
//   user=jack password=secret host=pg.example.com port=5432 dbname=mydb sslmode=verify-ca pool_max_conns=10
//
//   # Example URL
//   postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10
func Connect(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}
	return pool, pool.Ping(ctx)
}
