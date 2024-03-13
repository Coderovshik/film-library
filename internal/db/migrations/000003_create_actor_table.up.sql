CREATE TABLE IF NOT EXISTS actor(
    actor_id SERIAL PRIMARY KEY,
    actor_name VARCHAR NOT NULL,
    sex VARCHAR NOT NULL,
    birthday DATE NOT NULL
);