package actor

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Coderovshik/film-library/internal/db"
)

var (
	ErrActorNotExist = errors.New("actor does not exist")
	ErrEmptyUpdate   = errors.New("no updates to apply")
)

var _ ActorRepository = (*Repository)(nil)

type Repository struct {
	db db.DBTX
}

func NewRepository(db db.DBTX) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetActor(ctx context.Context, id int32) (*Actor, error) {
	const op = "actor.Repository.GetActor"

	const query = `
		SELECT a.actor_id, a.actor_name, a.sex, a.birthday,
			STRING_AGG (m.movie_name, ';') movie_list
		FROM actor a 
		INNER JOIN actor_in_movie am USING (actor_id)
		INNER JOIN movie m USING (movie_id)
		WHERE a.actor_id = $1
		GROUP BY a.actor_id`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("ERROR: failed to prepare query\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var a Actor
	var filmString string
	err = stmt.QueryRowContext(ctx, id).Scan(&a.ID, &a.Name, &a.Sex, &a.Birthday, &filmString)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("ERROR: actor with id=%d does not exist\n", id)
			return nil, fmt.Errorf("%s: %w", op, ErrActorNotExist)
		}

		log.Printf("ERROR: failed to execute query\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	a.Films = strings.Split(filmString, ";")

	return &a, nil
}

func (r *Repository) AddActor(ctx context.Context, a *Actor) (*Actor, error) {
	const op = "actor.Repository.AddActor"

	const query = `
		INSERT INTO actor(actor_name, sex, birthday)
		VALUES ($1, $2, $3) RETURNING actor_id`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("ERROR: failed to prepare query\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, a.Name, a.Sex, a.Birthday).Scan(&a.ID)
	if err != nil {
		log.Printf("ERROR: failed to execute query\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return a, nil
}

func (r *Repository) DeleteActor(ctx context.Context, id int32) error {
	const op = "actor.Repository.DeleteActor"

	const query = `DELETE FROM actor WHERE actor_id = $1`
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
		return fmt.Errorf("%s: %w", op, ErrActorNotExist)
	}

	return nil
}

func (r *Repository) UpdateActor(ctx context.Context, a *Actor) error {
	const op = "actor.Repository.UpdateActor"

	qo := ToQueryableObject(a)
	if qo.IsEmpty() {
		log.Print("ERROR: no updates to apply\n")
		return fmt.Errorf("%s: %w", op, ErrEmptyUpdate)
	}

	query := `UPDATE actor SET ` + qo.Args(2) + ` WHERE actor_id = $1`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("ERROR: failed to prepare query\n")
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, a.ID)
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
		log.Printf("ERROR: zero rows affected by deletion\n")
		return fmt.Errorf("%s: %w", op, ErrActorNotExist)
	}

	return nil
}

func (r *Repository) GetActors(ctx context.Context) ([]*Actor, error) {
	const op = "actor.Repository.GetActors"

	const query = `
		SELECT a.actor_id, a.actor_name, a.sex, a.birthday,
    		STRING_AGG (m.movie_name, ';') movie_list
		FROM actor a
		INNER JOIN actor_in_movie am USING (actor_id)
		INNER JOIN movie m USING (movie_id)
		GROUP BY a.actor_id`
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

	var actors []*Actor

	for rows.Next() {
		var a Actor
		var filmString string
		err := rows.Scan(&a.ID, &a.Name, &a.Sex, &a.Birthday, &filmString)
		if err != nil {
			log.Printf("ERROR: failed to execute query\n")
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		a.Films = strings.Split(filmString, ";")

		actors = append(actors, &a)
	}
	if err := rows.Err(); err != nil {
		log.Printf("ERROR: failed to execute query\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return actors, nil
}
