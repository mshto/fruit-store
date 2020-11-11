CREATE TABLE users_cart (
    product_id uuid REFERENCES products(id),
    user_id uuid REFERENCES users(id),
    amount INT,
    CONSTRAINT product_user_pkey PRIMARY KEY (product_id, user_id)
);