package service

import (
	"errors"
	"go-ecommerce/internal/dto"
	apperrors "go-ecommerce/internal/errors"
	"go-ecommerce/internal/models"
	"go-ecommerce/internal/repository"
	"strings"
)

type ProductService interface {
	Create(req dto.CreateProductRequest) (*dto.ProductResponse, error)
	GetByID(id uint) (*dto.ProductResponse, error)
	List(limit, offset int) ([]dto.ProductResponse, int64, error)
	Update(id uint, req dto.CreateProductRequest) (*dto.ProductResponse, error)
	Delete(id uint) error
}

type productService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productService{productRepo: productRepo}
}

func (p *productService) Create(req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	pr := models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		SKU:         req.SKU,
		IsActive:    req.IsActive,
	}

	_, err := p.productRepo.FindBySKU(pr.SKU)

	if err == nil {
		return nil, apperrors.ErrProductExists
	}

	if !errors.Is(err, apperrors.ErrProductNotFound) {
		return nil, apperrors.ErrInternal
	}

	if err = p.productRepo.Create(&pr); err != nil {

		errStr := err.Error()
		if strings.Contains(errStr, "Duplicate entry") || strings.Contains(errStr, "unique constraint") || strings.Contains(errStr, "UNIQUE constraint failed") {
			return nil, apperrors.ErrProductExists
		}

		return nil, apperrors.ErrInternal
	}

	response := dto.ToProductResponse(&pr)
	return &response, nil
}

func (p *productService) GetByID(id uint) (*dto.ProductResponse, error) {
	pr, err := p.productRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, apperrors.ErrProductNotFound) {
			return nil, apperrors.ErrProductNotFound
		}
		return nil, apperrors.ErrInternal
	}

	response := dto.ToProductResponse(pr)
	return &response, nil
}

func (p *productService) List(limit, offset int) ([]dto.ProductResponse, int64, error) {
	pr, err := p.productRepo.FindAll(limit, offset)
	if err != nil {
		return nil, 0, apperrors.ErrInternal
	}

	count, err := p.productRepo.Count()
	if err != nil {
		return nil, 0, apperrors.ErrInternal
	}

	res := make([]dto.ProductResponse, len(pr))
	for i := range pr {
		res[i] = dto.ToProductResponse(&pr[i])
	}
	return res, count, nil
}

func (p *productService) Update(id uint, req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	pr, err := p.productRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, apperrors.ErrProductNotFound) {
			return nil, apperrors.ErrProductNotFound
		}
		return nil, apperrors.ErrInternal
	}

	pr.Name = req.Name
	pr.Description = req.Description
	pr.Price = req.Price
	pr.Stock = req.Stock
	pr.IsActive = req.IsActive

	if err := p.productRepo.Update(pr); err != nil {
		return nil, apperrors.ErrInternal
	}

	response := dto.ToProductResponse(pr)
	return &response, nil
}

func (p *productService) Delete(id uint) error {
	_, err := p.productRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, apperrors.ErrProductNotFound) {
			return apperrors.ErrProductNotFound
		}
		return apperrors.ErrInternal
	}

	if err := p.productRepo.Delete(id); err != nil {
		return apperrors.ErrInternal
	}
	return nil
}
