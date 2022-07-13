-- name: CreateUser :one
INSERT INTO users (fname, lname, email, enc_passwd) 
VALUES (@fname::varchar(100), @lname::varchar(100), @email::varchar(100), crypt(@passwd::varchar, gen_salt('bf'))) 
RETURNING id, fname, lname, email, enc_passwd, created_at;

-- name: DeleteUserByID :execrows
DELETE FROM users WHERE id = $1;

-- name: GetUserByID :one
SELECT id, fname, lname, email, created_at FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT id, fname, lname, email, created_at FROM users
WHERE 
id > @id::INTEGER AND
(CASE WHEN @where_email::boolean 
    THEN  @email::VARCHAR(100) ILIKE email ELSE true END) AND
(CASE WHEN @where_fname::boolean 
    THEN  @fname::VARCHAR(100) ILIKE fname ELSE true END) AND
(CASE WHEN @where_lname::boolean 
    THEN  @lname::VARCHAR(100) ILIKE lname ELSE true END) 
ORDER BY id
LIMIT CASE WHEN @use_rlimit::boolean THEN @rlimit::INT END;

