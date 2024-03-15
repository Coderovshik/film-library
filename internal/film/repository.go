package film

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Coderovshik/film-library/internal/db"
	"github.com/lib/pq"
)

var (
	ErrFilmNotExist   = errors.New("actor does not exist")
	ErrEmptyUpdate    = errors.New("no updates to apply")
	ErrFilmActorExist = errors.New("given film and actor are already bound")
	ErrActorNotExist  = errors.New("actor with given id does not exist")
)

var _ FilmRepository = (*Repository)(nil)

type Repository struct {
	db db.DBTX
}

func NewRepository(db db.DBTX) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetFilm(ctx context.Context, id int32) (*Film, error) {
	const op = "film.Repository.GetFilm"

	const query = `
		SELECT m.movie_id, m.movie_name, m.movie_description, m.releasedate,
    		m.rating, STRING_AGG (a.actor_name, ';') movie_list
		FROM movie m
		LEFT JOIN actor_in_movie am USING (movie_id)
		LEFT JOIN actor a USING (actor_id)
		WHERE m.movie_id = $1
		GROUP BY m.movie_id`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("ERROR: failed to prepare query\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var f Film
	var actorString sql.NullString
	err = stmt.QueryRowContext(ctx, id).Scan(&f.ID, &f.Name, &f.Description, &f.ReleaseDate, &f.Rating, &actorString)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("ERROR: actor with id=%d does not exist\n", id)
			return nil, fmt.Errorf("%s: %w", op, ErrFilmNotExist)
		}

		log.Printf("ERROR: failed to execute query\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if len(actorString.String) != 0 {
		f.Actors = strings.Split(actorString.String, ";")
	}

	return &f, nil
}

func (r *Repository) AddFilm(ctx context.Context, f *Film) (*Film, error) {
	const op = "film.Repository.AddFilm"

	const query = `
		INSERT INTO movie(movie_name, movie_description, releasedate, rating)
		VALUES ($1, $2, $3, $4) RETURNING movie_id`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("ERROR: failed to prepare query\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, f.Name, f.Description, f.ReleaseDate, f.Rating).Scan(&f.ID)
	if err != nil {
		log.Printf("ERROR: failed to execute query\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return f, nil
}

func (r *Repository) AddFilmActors(ctx context.Context, fa *FilmActors) error {
	const op = "film.Repository.AddFilmActors"

	args, values := ToQueryableLists(fa, "($%d, $1)")
	query := `INSERT INTO actor_in_movie(actor_id, movie_id) VALUES ` + args
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("ERROR: failed to prepare query\n")
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, values...)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			if pgErr.Code.Name() == "unique_violation" {
				log.Printf("ERROR: one of the film-actor pairs already exists\n")
				return fmt.Errorf("%s: %w", op, ErrFilmActorExist)
			}
		}

		if errors.As(err, &pgErr) {
			if pgErr.Code.Name() == "foreign_key_violation" {
				log.Printf("ERROR: one of the film-actor pairs already exists\n")
				return fmt.Errorf("%s: %w", op, ErrActorNotExist)
			}
		}

		log.Printf("ERROR: failed to execute query\n")
		return fmt.Errorf("%s: %w", op, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		log.Printf("ERROR: failed to retrieve amount of rows affected by query\n")
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Printf("INFO: %d rows inserted\n", count)

	return nil
}

func (r *Repository) DeleteFilm(ctx context.Context, id int32) error {
	const op = "film.Repository.DeleteFilm"

	const query = `DELETE FROM movie WHERE movie_id = $1`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("ERROR: failed to prepare query\n")
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		log.Printf("ERROR: failed to execute query\n")
		return fmt.Errorf("%s: %w", op, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		log.Printf("ERROR: failed to retrieve amount of rows affected by query\n")
		return fmt.Errorf("%s: %w", op, err)
	}
	if count == 0 {
		log.Printf("ERROR: zero rows addected by deletion\n")
		return fmt.Errorf("%s: %w", op, ErrFilmNotExist)
	}

	return nil
}

func (r *Repository) UpdateFilm(ctx context.Context, f *Film) error {
	const op = "film.Repository.UpdateFilm"

	qo := ToQueryableObject(f)
	if qo.IsEmpty() {
		log.Print("ERROR: no updates to apply\n")
		return fmt.Errorf("%s: %w", op, ErrEmptyUpdate)
	}

	query := `UPDATE movie SET ` + qo.Args(1) +
		` WHERE movie_id = ` + fmt.Sprintf("$%d", qo.Len()+1)
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("ERROR: failed to prepare query\n")
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	values := qo.Values()
	values = append(values, f.ID)
	res, err := stmt.ExecContext(ctx, values...)
	if err != nil {
		log.Printf("ERROR: failed to execute query\n")
		return fmt.Errorf("%s: %w", op, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		log.Printf("ERROR: failed to retrieve amount of rows affected by query\n")
		return fmt.Errorf("%s: %w", op, err)
	}
	if count == 0 {
		log.Printf("ERROR: zero rows affected by update\n")
		return fmt.Errorf("%s: %w", op, ErrFilmNotExist)
	}

	return nil
}

func (r *Repository) GetFilms(ctx context.Context, q *Query) ([]*Film, error) {
	const op = "film.Repository.GetFilms"

	cons := ToQueryConditions(q)
	query := `
		SELECT m.movie_id, m.movie_name, m.movie_description, m.releasedate,
			m.rating, STRING_AGG (a.actor_name, ';') movie_list
		FROM movie m 
		LEFT JOIN actor_in_movie am USING (movie_id)
		LEFT JOIN actor a USING (actor_id) ` +
		cons[0] + " GROUP BY m.movie_id " +
		cons[1] + " " + cons[2]
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("ERROR: failed to prepare query\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		log.Printf("ERROR: failed to execute query\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var films []*Film

	for rows.Next() {
		var f Film
		var actorString sql.NullString
		err := rows.Scan(&f.ID, &f.Name, &f.Description, &f.ReleaseDate, &f.Rating, &actorString)
		if err != nil {
			log.Printf("ERROR: failed to execute query\n")
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		if len(actorString.String) != 0 {
			f.Actors = strings.Split(actorString.String, ";")
		}

		films = append(films, &f)
	}
	if err := rows.Err(); err != nil {
		log.Printf("ERROR: failed to execute query\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return films, nil
}

func (r *Repository) GetFilmActors(ctx context.Context, id int32) ([]*ActorShort, error) {
	const op = "film.Repository.GetFilmActors"

	const query = `
		SELECT a.actor_id, a.actor_name
		FROM actor a
		INNER JOIN actor_in_movie am USING (actor_id)
		WHERE movie_id = $1`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("ERROR: failed to prepare query\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		log.Printf("ERROR: failed to execute query\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var actors []*ActorShort
	for rows.Next() {
		var as ActorShort
		err := rows.Scan(&as.ID, &as.Name)
		if err != nil {
			log.Printf("ERROR: failed to execute query\n")
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		actors = append(actors, &as)
	}
	if err := rows.Err(); err != nil {
		log.Printf("ERROR: failed to execute query\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return actors, nil
}

func (r *Repository) DeleteFilmActors(ctx context.Context, fa *FilmActors) error {
	const op = "film.Repository.DeleteDilmActors"

	args, values := ToQueryableLists(fa, "$%d")
	var query = "DELETE FROM actor_in_movie WHERE movie_id = $1 AND actor_id IN (" + args + ")"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("ERROR: failed to prepare query\n")
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, values...)
	if err != nil {
		log.Printf("ERROR: failed to execute query\n")
		return fmt.Errorf("%s: %w", op, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		log.Printf("ERROR: failed to retrieve amount of rows affected by query\n")
		return fmt.Errorf("%s: %w", op, err)
	}
	if count == 0 {
		log.Printf("ERROR: zero rows affected by update\n")
		return fmt.Errorf("%s: %w", op, ErrFilmNotExist)
	}

	return nil
}
