CREATE TABLE discount (
    id VARCHAR(20) NOT NULL UNIQUE,
    rule VARCHAR(256),
    elements json NOT NULL,
    discount INT,
    PRIMARY KEY (id)
);

INSERT INTO discount
  (id, rule, elements, discount)
  VALUES ('luSjeDk3hd', 'more', '{"Oranges":1}', 30);

INSERT INTO discount
  (id, rule, elements, discount)
  VALUES ('9yRbSM4yd7', 'more', '{"Oranges":1}', 30);

INSERT INTO discount
  (id, rule, elements, discount)
  VALUES ('0fjotDPz3p', 'more', '{"Oranges":1}', 30);

INSERT INTO discount
  (id, rule, elements, discount)
  VALUES ('y0iq3P8E6T', 'more', '{"Oranges":1}', 30);

INSERT INTO discount
  (id, rule, elements, discount)
  VALUES ('EFeUYy6s6k', 'more', '{"Oranges":1}', 30);
