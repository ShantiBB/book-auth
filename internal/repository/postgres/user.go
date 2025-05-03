package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"auth/internal/entity"
)

func (r *Repository) CreateUser(ctx context.Context, u *entity.User) error {
	query := `INSERT INTO users (username, email, password_hash)
			  VALUES ($1, $2, $3) RETURNING id, role`

	var pgErr *pgconn.PgError
	err := r.db.QueryRow(ctx, query, u.Username, u.Email, u.PasswordHash).Scan(&u.ID, &u.Role)

	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return DuplicateError
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetUserCredentialsByUsername(ctx context.Context, u *entity.User) error {
	query := `SELECT id, password_hash, role FROM users WHERE username = $1`

	err := r.db.QueryRow(ctx, query, u.Username).Scan(&u.ID, &u.PasswordHash, &u.Role)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetUserByID(ctx context.Context, u *entity.User) error {
	query := `SELECT username, email, role, created_at, updated_at
			  FROM users WHERE id = $1`

	err := r.db.QueryRow(ctx, query, u.ID).Scan(
		&u.Username,
		&u.Email,
		&u.Role,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetAllUsers(ctx context.Context) ([]*entity.User, error) {
	query := `SELECT id, username, email, role
			  FROM users`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*entity.User

	for rows.Next() {
		var u entity.User
		if err = rows.Scan(&u.ID, &u.Username, &u.Email, &u.Role); err != nil {
			return nil, err
		}

		users = append(users, &u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Repository) UpdateUserByID(ctx context.Context, u *entity.User) error {
	query := `UPDATE users 
			  SET username = $1, email = $2, updated_at = NOW()
			  WHERE id = $3
			  RETURNING role`

	err := r.db.QueryRow(ctx, query, u.Username, u.Email, u.ID).Scan(&u.Role)
	if errors.Is(err, pgx.ErrNoRows) {
		return sql.ErrNoRows
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteUserByID(ctx context.Context, id int64) error {
	query := `DELETE FROM users 
       		  WHERE id = $1`

	res, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
