package dto

import "go-ecommerce/internal/models"

type AddToCartRequest struct {
	ProductID uint `json:"product_id" validate:"required"`
	Quantity  int  `json:"quantity" validate:"required,min=1"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" validate:"required,min=1"`
}

type CartItemResponse struct {
	ID          uint   `json:"id"`
	ProductID   uint   `json:"product_id"`
	ProductName string `json:"product_name"` // dari item.Product.Name
	Price       int64  `json:"price"`        // dari item.Product.Price
	Quantity    int    `json:"quantity"`
	Subtotal    int64  `json:"subtotal"` // dari item.Product.Price * item.Quantity
}

type CartResponse struct {
	ID         uint               `json:"id"`
	Items      []CartItemResponse `json:"items"`
	TotalPrice int64              `json:"total_price"` // sum dari subtotal
	TotalItems int                `json:"total_items"` // sum jumlah quantity item
}

func ToCartItemResponse(item *models.CartItem) CartItemResponse {
	return CartItemResponse{
		ID:          item.ID,
		ProductID:   item.ProductID,
		ProductName: item.Product.Name,
		Price:       item.Product.Price,
		Quantity:    item.Quantity,
		Subtotal:    item.Product.Price * int64(item.Quantity),
	}
}

func ToCartResponse(cart *models.Cart) CartResponse {
	responses := make([]CartItemResponse, 0)
	var total int64
	var totalQty int

	for _, item := range cart.Items {
		responses = append(responses, ToCartItemResponse(&item))
		total += item.Product.Price * int64(item.Quantity)
		totalQty += item.Quantity
	}

	return CartResponse{
		ID:         cart.ID,
		Items:      responses,
		TotalPrice: total,
		TotalItems: totalQty,
	}
}
