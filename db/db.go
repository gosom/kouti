package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

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
