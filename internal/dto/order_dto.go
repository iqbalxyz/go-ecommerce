package dto

import "go-ecommerce/internal/models"

type OrderItemResponse struct {
	ID          uint   `json:"id"`
	ProductID   uint   `json:"product_id"`
	ProductName string `json:"product_name"`
	Price       int64  `json:"price"`
	Quantity    int    `json:"quantity"`
	Subtotal    int64  `json:"subtotal"`
}

type OrderResponse struct {
	ID         uint                `json:"id"`
	UserID     uint                `json:"user_id"`
	TotalPrice int64               `json:"total_price"`
	Status     string              `json:"status"`
	CreatedAt  string              `json:"created_at"`
	OrderItems []OrderItemResponse `json:"items,omitempty"`
}

func ToOrderItemResponse(item *models.OrderItem) OrderItemResponse {
	return OrderItemResponse{
		ID:          item.ID,
		ProductID:   item.ProductID,
		ProductName: item.ProductName,
		Price:       item.Price,
		Quantity:    item.Quantity,
		Subtotal:    item.Subtotal,
	}
}

func ToOrderResponse(order *models.Order) OrderResponse {
	items := make([]OrderItemResponse, 0)
	if order.Items != nil {
		for _, item := range order.Items {
			items = append(items, ToOrderItemResponse(&item))
		}
	}

	if len(items) == 0 {
		return OrderResponse{
			ID:         order.ID,
			UserID:     order.UserID,
			TotalPrice: order.TotalPrice,
			Status:     string(order.Status),
			CreatedAt:  order.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}
	return OrderResponse{
		ID:         order.ID,
		UserID:     order.UserID,
		TotalPrice: order.TotalPrice,
		Status:     string(order.Status),
		CreatedAt:  order.CreatedAt.Format("2006-01-02 15:04:05"),
		OrderItems: items,
	}
}
