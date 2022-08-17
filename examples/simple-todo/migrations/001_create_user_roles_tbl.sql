CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users(
    id INT GENERATED ALWAYS AS IDENTITY NOT NULL,
	uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
	identity VARCHAR(100) NOT NULL UNIQUE,
    enc_passwd TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL default (NOW() AT TIME ZONE 'utc'),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL default (NOW() AT TIME ZONE 'utc'),
    PRIMARY KEY(id)
);

CREATE OR REPLACE FUNCTION trigger_set_users_timestamp() RETURNS TRIGGER AS $$
	BEGIN
		NEW.updated_at = (NOW() AT TIME ZONE 'utc');
  		RETURN NEW;
	END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_users_timestamp
	BEFORE UPDATE ON users
	FOR EACH ROW
	EXECUTE PROCEDURE trigger_set_users_timestamp();

CREATE TABLE roles(
	id INT GENERATED ALWAYS AS IDENTITY NOT NULL,
	name VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL default (NOW() AT TIME ZONE 'utc'),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL default (NOW() AT TIME ZONE 'utc'),
	PRIMARY KEY(id)
);

CREATE TABLE users_roles(
	user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	role_id INT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
	PRIMARY KEY(user_id, role_id)
);

---- create above / drop below ----

DROP TABLE IF EXISTS users_roles;
DROP TRIGGER IF EXISTS set_users_timestamp ON users;
DROP FUNCTION IF EXISTS set_users_timestamp;
DROP TABLE IF EXISTS roles;
DROP TABLE users;
