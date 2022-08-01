package um

const insertRolesQ = `
INSERT INTO roles(name) 
SELECT unnest($1::text[]) AS name 
ON CONFLICT DO NOTHING
RETURNING id, name, created_at, updated_at
`

const insertUserRolesQ = `
INSERT INTO users_roles(user_id, role_id) 
SELECT $1 AS user_id, unnest($2::int[]) as role_id
ON CONFLICT(user_id, role_id)
DO NOTHING
`

const selectRolesQ = `
SELECT id, name, created_at, updated_at
FROM roles
WHERE 
(CASE WHEN $1::boolean THEN id = ANY($2) ELSE true END) AND
(CASE WHEN $3::boolean THEN name = ANY($4) ELSE true END)`
