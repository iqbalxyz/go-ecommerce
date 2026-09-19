package service

import (
	"errors"
	"go-ecommerce/internal/dto"
	apperrors "go-ecommerce/internal/errors"
	"go-ecommerce/internal/models"
	"go-ecommerce/internal/repository"
)

type OrderService interface {
	Checkout(userID uint) (*dto.OrderResponse, error)
	GetByID(userID uint, orderID uint) (*dto.OrderResponse, error)
	List(userID uint, limit, offset int) ([]dto.OrderResponse, int64, error)
}

type orderService struct {
	orderRepo repository.OrderRepository
	cartRepo  repository.CartRepository
}

func NewOrderService(orderRepo repository.OrderRepository, cartRepo repository.CartRepository) OrderService {
	return &orderService{orderRepo: orderRepo, cartRepo: cartRepo}
}

func (o *orderService) Checkout(userID uint) (*dto.OrderResponse, error) {
	cart, err := o.cartRepo.FindByUserID(userID)
	if errors.Is(err, apperrors.ErrCartNotFound) {
		return nil, apperrors.ErrEmptyCart
	}

	if len(cart.Items) == 0 {
		return nil, apperrors.ErrEmptyCart
	}

	var totalPrice int64
	orderItems := make([]models.OrderItem, 0)
	stockUpdate := make(map[uint]int)

	for _, item := range cart.Items {
		if item.Product == nil {
			return nil, apperrors.ErrInternal
		}

		if !item.Product.IsActive {
			return nil, apperrors.ErrProductInactive
		}
		if item.Product.Stock < item.Quantity {
			return nil, apperrors.ErrInsufficientStock
		}

		subTotal := item.Product.Price * int64(item.Quantity)
		totalPrice += subTotal

		orderItems = append(orderItems, models.OrderItem{
			ProductID:   item.ProductID,
			ProductName: item.Product.Name,
			Price:       item.Product.Price,
			Quantity:    item.Quantity,
			Subtotal:    subTotal,
		})

		stockUpdate[item.ProductID] = item.Product.Stock - item.Quantity
	}

	order := &models.Order{
		UserID:     userID,
		TotalPrice: totalPrice,
		Status:     "pending",
	}

	err = o.orderRepo.Checkout(repository.CheckoutData{
		Order:        order,
		Items:        orderItems,
		StockUpdates: stockUpdate,
	}, cart.ID)

	createdOrder, err := o.orderRepo.FindByID(order.ID)
	if err != nil {
		return nil, apperrors.ErrInternal
	}

	response := dto.ToOrderResponse(createdOrder)
	return &response, nil
}

func (o *orderService) GetByID(userID uint, orderID uint) (*dto.OrderResponse, error) {
	findOrder, err := o.orderRepo.FindByID(orderID)
	if err != nil {
		if errors.Is(err, apperrors.ErrOrderNotFound) {
			return nil, apperrors.ErrOrderNotFound
		}
		return nil, apperrors.ErrInternal
	}
	if findOrder.UserID != userID {
		return nil, apperrors.ErrOrderNotFound
	}

	response := dto.ToOrderResponse(findOrder)
	return &response, nil
}

func (o *orderService) List(userID uint, limit int, offset int) ([]dto.OrderResponse, int64, error) {
	findOrder, err := o.orderRepo.FindByUserID(userID, limit, offset)
	if err != nil {
		if errors.Is(err, apperrors.ErrOrderNotFound) {
			return nil, 0, apperrors.ErrOrderNotFound
		}
		return nil, 0, apperrors.ErrInternal
	}
	count, err := o.orderRepo.CountByUserID(userID)
	if err != nil {
		return nil, 0, apperrors.ErrInternal
	}
	var responses []dto.OrderResponse
	for _, order := range findOrder {
		responses = append(responses, dto.ToOrderResponse(&order))
	}
	return responses, count, nil
}
