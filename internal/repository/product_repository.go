package repository

import (
	"errors"
	apperrors "go-ecommerce/internal/errors"
	"go-ecommerce/internal/models"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *models.Product) error
	FindByID(id uint) (*models.Product, error)
	FindBySKU(sku string) (*models.Product, error)
	FindAll(limit, offset int) ([]models.Product, error)
	Count() (int64, error)
	Update(product *models.Product) error
	Delete(id uint) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrProductNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindBySKU(sku string) (*models.Product, error) {
	var product models.Product
	err := r.db.Where("sku = ?", sku).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrProductNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindAll(limit int, offset int) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Limit(limit).Offset(offset).Order("id DESC").Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Product{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *productRepository) Update(product *models.Product) error {
	err := r.db.Model(&models.Product{}).Where("id = ?", product.ID).Updates(map[string]interface{}{
		"name":        product.Name,
		"description": product.Description,
		"price":       product.Price,
		"stock":       product.Stock,
		"is_active":   product.IsActive,
	}).Error

	if err != nil {
		return err
	}
	return nil
}

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}
