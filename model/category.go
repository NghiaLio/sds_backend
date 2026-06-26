package model

import (
	"errors"
	"strings"
)

// Category represents a product category.
type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Validate checks category inputs.
func (c *Category) Validate() error {
	c.Name = strings.TrimSpace(c.Name)
	if c.Name == "" {
		return errors.New("category name is required")
	}
	return nil
}
