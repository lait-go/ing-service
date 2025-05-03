CREATE TABLE IF NOT EXISTS person (
  id serial PRIMARY KEY,
  name varchar NOT NULL,
  phone varchar NOT NULL,
  tg_nickname varchar
);