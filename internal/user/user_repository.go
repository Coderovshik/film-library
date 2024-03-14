package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/lib/pq"
)

var (
	ErrUserExist    = errors.New("user already exists")
	ErrUserNotExist = errors.New("user does not exist")
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type Repository struct {
	db DBTX
}

func NewRepository(db DBTX) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateUser(ctx context.Context, user *User) (*User, error) {
	const op = "user.Repository.CreateUser"

	const query = "INSERT INTO users(user_name, passhash, is_admin) VALUES ($1, $2, $3) RETURNING user_id"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("ERROR: failed to prepare query\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, user.Username, user.Passhash, user.IsAdmin).Scan(&user.ID)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			if pgErr.Code.Name() == "unique_violation" {
				log.Printf("ERROR: user %s already exists\n", user.Username)
				return nil, fmt.Errorf("%s: %w", op, ErrUserExist)
			}
		}

		log.Printf("ERROR: failed to execute query\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	const op = "user.Repository.GetUser"

	const query = "SELECT user_id, user_name, passhash, is_admin FROM users WHERE user_name = $1"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("ERROR: failed to prepare query\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	u := User{}
	err = stmt.QueryRowContext(ctx, username).Scan(&u.ID, &u.Username, &u.Passhash, &u.IsAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("ERROR: user %s does not exist\n", username)
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotExist)
		}

		log.Printf("ERROR: failed to execute query\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &u, nil
}
