package service

import (
	"context"
	"errors"
	"fmt"
	"mobile-crud-backend/model"
	"mobile-crud-backend/repository"
)

// ProductService handles business logic for products.
type ProductService interface {
	CreateProduct(ctx context.Context, product *model.Product) error
	GetProductByID(ctx context.Context, id int64) (*model.Product, error)
	GetAllProducts(ctx context.Context, categoryID int64, keyword string, page, limit int) ([]model.Product, int, error)
	UpdateProduct(ctx context.Context, id int64, product *model.Product) error
	DeleteProduct(ctx context.Context, id int64) error
}

type productServiceImpl struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
}

// NewProductService creates a new ProductService instance.
func NewProductService(productRepo repository.ProductRepository, categoryRepo repository.CategoryRepository) ProductService {
	return &productServiceImpl{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *productServiceImpl) CreateProduct(ctx context.Context, product *model.Product) error {
	if err := product.Validate(); err != nil {
		return err
	}

	// 1. Verify Category exists
	cat, err := s.categoryRepo.GetByID(ctx, product.CategoryID)
	if err != nil {
		return err
	}
	if cat == nil {
		return fmt.Errorf("category_id %d does not exist", product.CategoryID)
	}

	// 2. Verify Code uniqueness
	existing, err := s.productRepo.GetByCode(ctx, product.Code)
	if err != nil {
		return err
	}
	if existing != nil {
		return fmt.Errorf("product code '%s' already exists", product.Code)
	}

	return s.productRepo.Create(ctx, product)
}

func (s *productServiceImpl) GetProductByID(ctx context.Context, id int64) (*model.Product, error) {
	p, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, errors.New("product not found")
	}
	return p, nil
}

func (s *productServiceImpl) GetAllProducts(ctx context.Context, categoryID int64, keyword string, page, limit int) ([]model.Product, int, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	return s.productRepo.GetAll(ctx, categoryID, keyword, page, limit)
}

func (s *productServiceImpl) UpdateProduct(ctx context.Context, id int64, product *model.Product) error {
	existing, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("product not found")
	}

	// 1. Validate incoming values
	product.ID = id
	if err := product.Validate(); err != nil {
		return err
	}

	// 2. Verify Category exists
	cat, err := s.categoryRepo.GetByID(ctx, product.CategoryID)
	if err != nil {
		return err
	}
	if cat == nil {
		return fmt.Errorf("category_id %d does not exist", product.CategoryID)
	}

	// 3. Verify Code uniqueness
	codeOwner, err := s.productRepo.GetByCode(ctx, product.Code)
	if err != nil {
		return err
	}
	if codeOwner != nil && codeOwner.ID != id {
		return fmt.Errorf("product code '%s' is already in use by another product", product.Code)
	}

	// 4. Update fields
	existing.Name = product.Name
	existing.Code = product.Code
	existing.Price = product.Price
	existing.Stock = product.Stock
	existing.CategoryID = product.CategoryID
	existing.Description = product.Description
	existing.Image = product.Image

	if err := s.productRepo.Update(ctx, existing); err != nil {
		return err
	}

	*product = *existing
	return nil
}

func (s *productServiceImpl) DeleteProduct(ctx context.Context, id int64) error {
	existing, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("product not found")
	}
	return s.productRepo.Delete(ctx, id)
}
