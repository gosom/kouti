package um

const setTimestampFnQ = `
CREATE OR REPLACE FUNCTION trigger_set_users_timestamp() RETURNS TRIGGER AS $$
	BEGIN
		NEW.updated_at = (NOW() AT TIME ZONE 'utc');
  		RETURN NEW;
	END;
$$ LANGUAGE plpgsql`

const createUsrTblQ = `
CREATE TABLE IF NOT EXISTS users(
    id INT GENERATED ALWAYS AS IDENTITY NOT NULL,
	uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
	identity VARCHAR(100) NOT NULL UNIQUE,
    enc_passwd TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL default (NOW() AT TIME ZONE 'utc'),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL default (NOW() AT TIME ZONE 'utc'),
    PRIMARY KEY(id)
)`

const usersDropTsTrigger = `
DROP TRIGGER IF EXISTS set_users_timestamp ON users
`

const usersSetTsTrigger = `
CREATE TRIGGER set_users_timestamp
	BEFORE UPDATE ON users
	FOR EACH ROW
	EXECUTE PROCEDURE trigger_set_users_timestamp();
`

const createRoleTblQ = `
CREATE TABLE IF NOT EXISTS roles(
	id INT GENERATED ALWAYS AS IDENTITY NOT NULL,
	name VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL default (NOW() AT TIME ZONE 'utc'),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL default (NOW() AT TIME ZONE 'utc'),
	PRIMARY KEY(id)
)`

const createUserRolesTblQ = `
CREATE TABLE IF NOT EXISTS users_roles(
	user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	role_id INT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
	PRIMARY KEY(user_id, role_id)
)
`
