package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// UserRepository handles user data persistence.
type UserRepository interface {
	Create(ctx context.Context, username, hashedPassword string) error
	GetByUsername(ctx context.Context, username string) (int64, string, error)
}

type sqliteUserRepository struct {
	db *sql.DB
}

// NewSQLiteUserRepository creates a new SQLite user repository instance.
func NewSQLiteUserRepository(db *sql.DB) UserRepository {
	return &sqliteUserRepository{db: db}
}

func (r *sqliteUserRepository) Create(ctx context.Context, username, hashedPassword string) error {
	query := `
		INSERT INTO users (username, password, created_at)
		VALUES (?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query, username, hashedPassword, time.Now())
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

func (r *sqliteUserRepository) GetByUsername(ctx context.Context, username string) (int64, string, error) {
	query := `
		SELECT id, password
		FROM users
		WHERE username = ?
	`
	var id int64
	var passwordHash string
	err := r.db.QueryRowContext(ctx, query, username).Scan(&id, &passwordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, "", nil // user not found, return 0 ID and empty hash
		}
		return 0, "", fmt.Errorf("failed to retrieve user: %w", err)
	}
	return id, passwordHash, nil
}
