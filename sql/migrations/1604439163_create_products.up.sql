CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE products (
    id uuid DEFAULT uuid_generate_v1() NOT NULL,
    name character varying(255) NOT NULL,
    price float(2),
    created_at date,
    PRIMARY KEY (id)
);

INSERT INTO products
  (id,  name, price, created_at)
  VALUES (uuid_generate_v1(), 'Apples', 1.72, TRANSACTION_TIMESTAMP());

INSERT INTO products
  (id,  name, price, created_at)
  VALUES (uuid_generate_v1(), 'Bananas', 2.34, TRANSACTION_TIMESTAMP());

INSERT INTO products
  (id,  name, price, created_at)
  VALUES (uuid_generate_v1(), 'Pears', 0.13, TRANSACTION_TIMESTAMP());

INSERT INTO products
  (id,  name, price, created_at)
  VALUES (uuid_generate_v1(), 'Oranges', 1.65, TRANSACTION_TIMESTAMP());

  