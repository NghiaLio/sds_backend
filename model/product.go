package model

import (
	"errors"
	"strings"
	"time"
)

// Product represents an item in the catalog.
type Product struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CategoryID  int64     `json:"category_id"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Validate checks product fields.
func (p *Product) Validate() error {
	p.Name = strings.TrimSpace(p.Name)
	if p.Name == "" {
		return errors.New("product name is required")
	}

	p.Code = strings.TrimSpace(p.Code)
	if p.Code == "" {
		return errors.New("product code is required")
	}

	if p.Price < 0 {
		return errors.New("price cannot be negative")
	}

	if p.Stock < 0 {
		return errors.New("stock cannot be negative")
	}

	if p.CategoryID <= 0 {
		return errors.New("valid category_id is required")
	}

	return nil
}
