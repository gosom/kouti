package um

const usersCreateQ = `
INSERT INTO users(identity, enc_passwd) 
VALUES ($1, crypt($2, gen_salt('bf'))) 
RETURNING id, uid, identity, enc_passwd, created_at, updated_at
`

const usersGetQ = `
SELECT id, uid, identity, enc_passwd, created_at, updated_at
FROM users 
WHERE
(CASE WHEN $1::boolean THEN uid = $2 ELSE true END) AND
(CASE WHEN $3::boolean THEN id = $4 ELSE true END) AND
(CASE WHEN $5::boolean THEN identity ILIKE $6 ELSE true END)
`

const usersDeleteQ = `
DELETE
FROM users 
WHERE
(CASE WHEN $1::boolean THEN uid = $2::UUID ELSE true END) AND
(CASE WHEN $3::boolean THEN id = $4::INTEGER ELSE true END) AND
(CASE WHEN $5::boolean THEN identity ILIKE $6::VARCHAR(100) ELSE true END)
`

const usersSelectQ = `
SELECT id, uid, identity, enc_passwd, created_at, updated_at
FROM users 
WHERE
id > $1
ORDER BY id
LIMIT CASE WHEN $2::boolean THEN $3::INT END;
`
