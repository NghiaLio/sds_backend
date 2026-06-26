package service

import (
	"context"
	"errors"
	"mobile-crud-backend/model"
	"mobile-crud-backend/repository"
)

// CategoryService handles business logic for categories.
type CategoryService interface {
	CreateCategory(ctx context.Context, category *model.Category) error
	GetCategoryByID(ctx context.Context, id int64) (*model.Category, error)
	GetAllCategories(ctx context.Context) ([]model.Category, error)
	UpdateCategory(ctx context.Context, id int64, category *model.Category) error
	DeleteCategory(ctx context.Context, id int64) error
}

type categoryServiceImpl struct {
	repo repository.CategoryRepository
}

// NewCategoryService creates a new CategoryService instance.
func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryServiceImpl{repo: repo}
}

func (s *categoryServiceImpl) CreateCategory(ctx context.Context, category *model.Category) error {
	if err := category.Validate(); err != nil {
		return err
	}
	return s.repo.Create(ctx, category)
}

func (s *categoryServiceImpl) GetCategoryByID(ctx context.Context, id int64) (*model.Category, error) {
	cat, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, errors.New("category not found")
	}
	return cat, nil
}

func (s *categoryServiceImpl) GetAllCategories(ctx context.Context) ([]model.Category, error) {
	return s.repo.GetAll(ctx)
}

func (s *categoryServiceImpl) UpdateCategory(ctx context.Context, id int64, category *model.Category) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("category not found")
	}

	existing.Name = category.Name
	if err := existing.Validate(); err != nil {
		return err
	}

	if err := s.repo.Update(ctx, existing); err != nil {
		return err
	}

	*category = *existing
	return nil
}

func (s *categoryServiceImpl) DeleteCategory(ctx context.Context, id int64) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("category not found")
	}
	return s.repo.Delete(ctx, id)
}
