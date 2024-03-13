CREATE TABLE actor_in_movie(
    role_id  SERIAL PRIMARY KEY,
    actor_id INT NOT NULL REFERENCES actor(actor_id) ON DELETE CASCADE,
    movie_id INT NOT NULL REFERENCES movie(movie_id) ON DELETE CASCADE
);