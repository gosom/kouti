package db

import (
	"context"

	"github.com/gosom/kouti/dbdriver"
	"github.com/gosom/kouti/examples/todo/orm"
)

type DB struct {
	dbdriver.DB
	*orm.Queries
}

func New(ctx context.Context, cfg dbdriver.Config) (DB, error) {
	ans := DB{}
	d, err := dbdriver.New(ctx, cfg)
	if err != nil {
		return ans, err
	}
	ans.DB = d
	ans.Queries = orm.New(ans.DB.RawConn())
	return ans, nil
}

const userFtsQ = `SELECT 
     users.id,
     users.fname,
     users.lname,
     users.email,
     users.created_at
 FROM 
     users, 
     to_tsvector(users.email || users.fname || users.lname) document,
     plainto_tsquery($1) query,
     NULLIF(ts_rank(to_tsvector(users.email), query), 0) rank_email,
     NULLIF(ts_rank(to_tsvector(users.fname), query), 0) rank_fname,
     NULLIF(ts_rank(to_tsvector(users.lname), query), 0) rank_lname,
     SIMILARITY($1::TEXT, users.email || users.fname || users.lname) similarity
WHERE 
id > $2 AND (query @@ document OR similarity > 0)
ORDER BY rank_email, rank_lname, rank_fname, similarity DESC NULLS LAST
LIMIT $3
`

type UserFtsParams struct {
	ID     int
	Phrase string
	Limit  int
}

func (o DB) SearchUsers(ctx context.Context, p UserFtsParams) ([]orm.ListUsersRow, error) {
	rows, err := o.RawConn().Query(ctx, userFtsQ, p.Phrase, p.ID, p.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []orm.ListUsersRow
	for rows.Next() {
		var i orm.ListUsersRow
		if err := rows.Scan(
			&i.ID,
			&i.Fname,
			&i.Lname,
			&i.Email,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

// This complexity is done in order to apply bcrypt always
// this way we make enumeration attacks difficult
// thanks https://daniel.fone.net.nz/blog/2020/09/09/timing-safe-bcrypt-authentication-in-postgresql/
const userLoginQ = `
WITH
target_user AS (
  SELECT id, enc_passwd
  FROM (
    SELECT id, enc_passwd from users WHERE email = $1
    UNION ALL
    SELECT NULL, gen_salt('bf')
  ) users
  LIMIT 1
),
valid_user AS (
  SELECT id FROM target_user WHERE enc_passwd = crypt($2, enc_passwd)
)
SELECT id, fname, lname, email, created_at FROM users NATURAL JOIN valid_user LIMIT 1
`

type UserLoginParams struct {
	Email  string
	Passwd string
}

func (o DB) GetUserByEmailPasswd(ctx context.Context, p UserLoginParams) (orm.User, error) {
	row := o.RawConn().QueryRow(ctx, userLoginQ, p.Email, p.Passwd)
	var ans orm.User
	err := row.Scan(
		&ans.ID,
		&ans.Fname,
		&ans.Lname,
		&ans.Email,
		&ans.CreatedAt,
	)
	return ans, err
}
