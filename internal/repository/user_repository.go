package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/andrebarone77/cardiaflow-api/internal/domain"
	"github.com/lib/pq"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Save(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (name, email, password_hash)
		VALUES ($1, $2, $3)
		ON CONFLICT (email)
		DO UPDATE SET
			name = $1,
			password_hash = $3,
			updated_at = now()
		RETURNING id, created_at, updated_at;
	`

	err := r.db.QueryRowContext(ctx,
		query,
		user.Name,
		user.Email,
		user.PasswordHash,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error creating user: %v", err)
		pqErr, ok := err.(*pq.Error)

		if ok && pqErr.Code == "23505" && pqErr.Constraint == "users_email_key" {
			return domain.ErrEmailAlreadyExists
		}

		return err
	}

	return nil

}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	query := `
		SELECT id, name, email, password_hash
		FROM users
		WHERE email = $1
	`
	err := r.db.QueryRowContext(ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetById(ctx context.Context, id string) (*domain.User, error) {
	user := &domain.User{}
	query := `
		SELECT id, name, email, password_hash
		FROM users
		WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx,
		query,
		id,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrUserNotFound
	}

	return nil

}
