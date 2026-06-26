package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mobile-crud-backend/model"
	"time"
)

// ProductRepository handles Product CRUD in the database.
type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) error
	GetByID(ctx context.Context, id int64) (*model.Product, error)
	GetByCode(ctx context.Context, code string) (*model.Product, error)
	GetAll(ctx context.Context, categoryID int64, keyword string, page, limit int) ([]model.Product, int, error)
	Update(ctx context.Context, product *model.Product) error
	Delete(ctx context.Context, id int64) error
}

type sqliteProductRepository struct {
	db *sql.DB
}

// NewSQLiteProductRepository creates a new SQLite product repository.
func NewSQLiteProductRepository(db *sql.DB) ProductRepository {
	return &sqliteProductRepository{db: db}
}

func (r *sqliteProductRepository) Create(ctx context.Context, product *model.Product) error {
	query := `
		INSERT INTO products (name, code, price, stock, category_id, description, image, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now

	res, err := r.db.ExecContext(ctx, query,
		product.Name,
		product.Code,
		product.Price,
		product.Stock,
		product.CategoryID,
		product.Description,
		product.Image,
		product.CreatedAt,
		product.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert product: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to retrieve last insert id: %w", err)
	}

	product.ID = id
	return nil
}

func (r *sqliteProductRepository) GetByID(ctx context.Context, id int64) (*model.Product, error) {
	query := `
		SELECT id, name, code, price, stock, category_id, description, image, created_at, updated_at
		FROM products
		WHERE id = ?
	`
	var p model.Product
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Code,
		&p.Price,
		&p.Stock,
		&p.CategoryID,
		&p.Description,
		&p.Image,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get product by id: %w", err)
	}
	return &p, nil
}

func (r *sqliteProductRepository) GetByCode(ctx context.Context, code string) (*model.Product, error) {
	query := `
		SELECT id, name, code, price, stock, category_id, description, image, created_at, updated_at
		FROM products
		WHERE code = ?
	`
	var p model.Product
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&p.ID,
		&p.Name,
		&p.Code,
		&p.Price,
		&p.Stock,
		&p.CategoryID,
		&p.Description,
		&p.Image,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get product by code: %w", err)
	}
	return &p, nil
}

func (r *sqliteProductRepository) GetAll(ctx context.Context, categoryID int64, keyword string, page, limit int) ([]model.Product, int, error) {
	// Base queries
	query := `SELECT id, name, code, price, stock, category_id, description, image, created_at, updated_at FROM products WHERE 1=1`
	countQuery := `SELECT COUNT(1) FROM products WHERE 1=1`
	var args []interface{}

	if categoryID > 0 {
		query += " AND category_id = ?"
		countQuery += " AND category_id = ?"
		args = append(args, categoryID)
	}

	if keyword != "" {
		keywordSearch := "%" + keyword + "%"
		query += " AND (name LIKE ? OR code LIKE ? OR description LIKE ?)"
		countQuery += " AND (name LIKE ? OR code LIKE ? OR description LIKE ?)"
		args = append(args, keywordSearch, keywordSearch, keywordSearch)
	}

	// 1. Get total count
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	// 2. Add sorting, limit and offset
	query += " ORDER BY created_at DESC"
	if limit > 0 {
		query += " LIMIT ?"
		offset := (page - 1) * limit
		query += " OFFSET ?"
		args = append(args, limit, offset)
	}

	// 3. Query products
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()

	products := []model.Product{}
	for rows.Next() {
		var p model.Product
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Code,
			&p.Price,
			&p.Stock,
			&p.CategoryID,
			&p.Description,
			&p.Image,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("row error: %w", err)
	}

	return products, total, nil
}

func (r *sqliteProductRepository) Update(ctx context.Context, product *model.Product) error {
	query := `
		UPDATE products
		SET name = ?, code = ?, price = ?, stock = ?, category_id = ?, description = ?, image = ?, updated_at = ?
		WHERE id = ?
	`
	product.UpdatedAt = time.Now()

	res, err := r.db.ExecContext(ctx, query,
		product.Name,
		product.Code,
		product.Price,
		product.Stock,
		product.CategoryID,
		product.Description,
		product.Image,
		product.UpdatedAt,
		product.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with id %d not found", product.ID)
	}

	return nil
}

func (r *sqliteProductRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM products WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with id %d not found", id)
	}

	return nil
}
