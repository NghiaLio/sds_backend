package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mobile-crud-backend/model"
)

// CategoryRepository handles Category CRUD in the database.
type CategoryRepository interface {
	Create(ctx context.Context, category *model.Category) error
	GetByID(ctx context.Context, id int64) (*model.Category, error)
	GetAll(ctx context.Context) ([]model.Category, error)
	Update(ctx context.Context, category *model.Category) error
	Delete(ctx context.Context, id int64) error
}

type sqliteCategoryRepository struct {
	db *sql.DB
}

// NewSQLiteCategoryRepository creates a new SQLite category repository.
func NewSQLiteCategoryRepository(db *sql.DB) CategoryRepository {
	return &sqliteCategoryRepository{db: db}
}

func (r *sqliteCategoryRepository) Create(ctx context.Context, category *model.Category) error {
	query := `INSERT INTO categories (name) VALUES (?)`
	res, err := r.db.ExecContext(ctx, query, category.Name)
	if err != nil {
		return fmt.Errorf("failed to insert category: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to retrieve last insert id: %w", err)
	}

	category.ID = id
	return nil
}

func (r *sqliteCategoryRepository) GetByID(ctx context.Context, id int64) (*model.Category, error) {
	query := `SELECT id, name FROM categories WHERE id = ?`
	var category model.Category
	err := r.db.QueryRowContext(ctx, query, id).Scan(&category.ID, &category.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	return &category, nil
}

func (r *sqliteCategoryRepository) GetAll(ctx context.Context) ([]model.Category, error) {
	query := `SELECT id, name FROM categories ORDER BY id ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	categories := []model.Category{}
	for rows.Next() {
		var category model.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row error: %w", err)
	}

	return categories, nil
}

func (r *sqliteCategoryRepository) Update(ctx context.Context, category *model.Category) error {
	query := `UPDATE categories SET name = ? WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, category.Name, category.ID)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("category with id %d not found", category.ID)
	}

	return nil
}

func (r *sqliteCategoryRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM categories WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("category with id %d not found", id)
	}

	return nil
}
