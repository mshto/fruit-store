CREATE TABLE users (
    id uuid DEFAULT uuid_generate_v1(),
    username VARCHAR(34) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    PRIMARY KEY (id)
);