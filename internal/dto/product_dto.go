package dto

import (
	"go-ecommerce/internal/models"
	"time"
)

type CreateProductRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=200"`
	Description string `json:"description" validate:"max=5000"`
	Price       int64  `json:"price" validate:"required,min=0"`
	Stock       int    `json:"stock" validate:"min=0"`
	SKU         string `json:"sku" validate:"required,min=3,max=50"`
	IsActive    bool   `json:"is_active"`
}

type UpdateProductRequest struct {
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Price       *int64     `json:"price"`
	Stock       *int       `json:"stock"`
	IsActive    *bool      `json:"is_active"`
	CreatedAt   *time.Time `json:"created_at"`
}

type ProductResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	Stock       int       `json:"stock"`
	SKU         string    `json:"sku"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

func ToProductResponse(p *models.Product) ProductResponse {
	//Fungsi konversi.
	return ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		SKU:         p.SKU,
		IsActive:    p.IsActive,
		CreatedAt:   p.CreatedAt,
	}

}

func ToProductResponses(products []models.Product) []ProductResponse {
	//Konversi slice (untuk list). Loop dan konversi tiap item.
	var responses []ProductResponse
	for _, product := range products {
		responses = append(responses, ToProductResponse(&product))
	}
	return responses
}
