CREATE TABLE actor_in_movie(
    actor_id INT NOT NULL REFERENCES actor(actor_id) ON DELETE CASCADE,
    movie_id INT NOT NULL REFERENCES movie(movie_id) ON DELETE CASCADE,
    PRIMARY KEY (actor_id, movie_id)
);