CREATE TABLE todos(
    id UUID NOT NULL,
    content VARCHAR(140) NOT NULL,
    is_completed BOOLEAN NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    PRIMARY KEY(id)
);

---- create above / drop below ----

DROP TABLE todos;
