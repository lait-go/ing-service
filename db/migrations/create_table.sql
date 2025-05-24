CREATE TABLE IF NOT EXISTS person (
  id serial PRIMARY KEY,
  profession varchar NOT NULL,
  name varchar NOT NULL,
  phone varchar NOT NULL,
  tg varchar
);