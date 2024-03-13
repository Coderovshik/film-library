CREATE TABLE IF NOT EXISTS movie(
    movie_id SERIAL PRIMARY KEY,
    movie_name VARCHAR NOT NULL,
    movie_description VARCHAR NOT NULL,
    releasedate DATE NOT NULL,
    rating INTEGER NOT NULL
);