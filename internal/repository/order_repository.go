package repository

import (
	"errors"
	apperrors "go-ecommerce/internal/errors"
	"go-ecommerce/internal/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateWithItems(order *models.Order, items []models.OrderItem) error
	FindByID(orderID uint) (*models.Order, error)
	FindByUserID(userID uint, limit, offset int) ([]models.Order, error)
	CountByUserID(userID uint) (int64, error)
	Checkout(data CheckoutData, clearCartID uint) error
}

type orderRepository struct {
	db *gorm.DB
}

type CheckoutData struct {
	Order        *models.Order
	Items        []models.OrderItem
	StockUpdates map[uint]int
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (o *orderRepository) CreateWithItems(order *models.Order, items []models.OrderItem) error {
	result := o.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		for i := range items {
			items[i].OrderID = order.ID
		}

		if err := tx.Create(&items).Error; err != nil {
			return err
		}

		return nil
	})
	return result
}

func (o *orderRepository) FindByID(orderID uint) (*models.Order, error) {
	var order models.Order
	err := o.db.Preload("Items").First(&order, orderID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrOrderNotFound
		}
		return nil, err
	}
	return &order, nil
}

func (o *orderRepository) FindByUserID(userID uint, limit int, offset int) ([]models.Order, error) {
	var orders []models.Order
	err := o.db.Where("user_id = ?", userID).
		Order("id DESC").
		Limit(limit).Offset(offset).
		Preload("Items").
		Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (o *orderRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := o.db.Model(&models.Order{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *orderRepository) Checkout(data CheckoutData, clearCartID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(data.Order).Error; err != nil {
			return err
		}
		for i := range data.Items {
			data.Items[i].OrderID = data.Order.ID
		}
		if err := tx.Create(&data.Items).Error; err != nil {
			return err
		}
		for productID, newStock := range data.StockUpdates {
			if err := tx.Model(&models.Product{}).
				Where("id = ?", productID).
				Update("stock", newStock).Error; err != nil {
				return err
			}
		}
		if clearCartID > 0 {
			if err := tx.Where("cart_id = ?", clearCartID).
				Delete(&models.CartItem{}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
