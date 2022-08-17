CREATE TABLE todos(
    id INT GENERATED ALWAYS AS IDENTITY NOT NULL,
	uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    content TEXT NOT NULL,
    user_id INT NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT false,
    PRIMARY KEY(id),
    CONSTRAINT fk_user 
        FOREIGN KEY(user_id) 
	    REFERENCES users(id)
        ON DELETE CASCADE
);

---- create above / drop below ----
DROP TABLE todos;
