CREATE TABLE IF NOT EXISTS users(
    user_id SERIAL PRIMARY KEY,
    user_name VARCHAR UNIQUE NOT NULL,
    passhash VARCHAR NOT NULL,
    is_admin BOOLEAN NOT NULL
);