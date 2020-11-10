CREATE TABLE users_cart (
    id uuid DEFAULT uuid_generate_v1(),
    user_id uuid REFERENCES users(id),
    count int,
    name character varying(255) NOT NULL,
    PRIMARY KEY (id)
);
