package repository

import (
	"errors"
	apperrors "go-ecommerce/internal/errors"
	"go-ecommerce/internal/models"

	"gorm.io/gorm"
)

type CartRepository interface {
	FindByUserID(userID uint) (*models.Cart, error)
	Create(cart *models.Cart) error

	FindItemByID(itemID uint) (*models.CartItem, error)
	FindItemByCartAndProduct(cartID, productID uint) (*models.CartItem, error)
	CreateItem(item *models.CartItem) error
	UpdateItem(item *models.CartItem) error
	DeleteItem(itemID uint) error
	ClearItems(cartID uint) error
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (c *cartRepository) FindByUserID(userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := c.db.Preload("Items.Product").Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrCartNotFound
		}
		return nil, err
	}
	return &cart, nil
}

func (c *cartRepository) Create(cart *models.Cart) error {
	return c.db.Create(cart).Error
}

func (c *cartRepository) FindItemByID(itemID uint) (*models.CartItem, error) {
	var item models.CartItem
	err := c.db.Preload("Product").First(&item, itemID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrCartItemNotFound
		}
		return nil, err
	}
	return &item, nil
}

func (c *cartRepository) FindItemByCartAndProduct(cartID uint, productID uint) (*models.CartItem, error) {
	var item models.CartItem
	err := c.db.Preload("Product").Where("cart_id = ? AND product_id = ?", cartID, productID).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrCartItemNotFound
		}
		return nil, err
	}
	return &item, nil
}

func (c *cartRepository) CreateItem(item *models.CartItem) error {
	// insert item
	err := c.db.Create(item).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *cartRepository) UpdateItem(item *models.CartItem) error {
	err := c.db.Model(&models.CartItem{}).Where("id = ?", item.ID).Updates(map[string]interface{}{
		"quantity": item.Quantity,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *cartRepository) DeleteItem(itemID uint) error {
	result := c.db.Where("id = ?", itemID).Delete(&models.CartItem{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return apperrors.ErrCartItemNotFound
	}

	return nil
}

func (c *cartRepository) ClearItems(cartID uint) error {
	return c.db.Where("cart_id = ?", cartID).Delete(&models.CartItem{}).Error
}
