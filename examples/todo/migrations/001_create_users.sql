CREATE EXTENSION pgcrypto;

CREATE TABLE users(
    id INT GENERATED ALWAYS AS IDENTITY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE default (NOW() AT TIME ZONE 'utc')
);

---- create above / drop below ----

drop table users;
