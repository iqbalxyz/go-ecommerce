package service

import (
	"errors"
	"go-ecommerce/internal/dto"
	apperrors "go-ecommerce/internal/errors"
	"go-ecommerce/internal/models"
	"go-ecommerce/internal/repository"
)

type CartService interface {
	GetCart(userID uint) (*dto.CartResponse, error)
	AddItem(userID uint, req dto.AddToCartRequest) (*dto.CartResponse, error)
	UpdateItem(userID uint, itemID uint, req dto.UpdateCartItemRequest) (*dto.CartResponse, error)
	RemoveItem(userID uint, itemID uint) (*dto.CartResponse, error)
}

type cartService struct {
	productRepo repository.ProductRepository
	cartRepo    repository.CartRepository
}

func NewCartService(cartRepo repository.CartRepository, productRepo repository.ProductRepository) CartService {
	return &cartService{cartRepo: cartRepo, productRepo: productRepo}
}

func (c *cartService) GetCart(userID uint) (*dto.CartResponse, error) {
	cart, err := c.cartRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrCartNotFound) {
			cart = &models.Cart{UserID: userID}
			if err := c.cartRepo.Create(cart); err != nil {
				return nil, apperrors.ErrInternal
			}
		} else {
			return nil, apperrors.ErrInternal
		}
	}
	response := dto.ToCartResponse(cart)
	return &response, nil
}

func (c *cartService) AddItem(userID uint, req dto.AddToCartRequest) (*dto.CartResponse, error) {
	cart, err := c.cartRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrCartNotFound) {
			cart = &models.Cart{
				UserID: userID,
			}
			if err := c.cartRepo.Create(cart); err != nil {
				return nil, apperrors.ErrInternal
			}
		} else {
			return nil, apperrors.ErrInternal
		}
	}

	product, err := c.productRepo.FindByID(req.ProductID)
	if err != nil {
		if errors.Is(err, apperrors.ErrProductNotFound) {
			return nil, apperrors.ErrProductNotFound
		}
		return nil, apperrors.ErrInternal
	}

	if !product.IsActive {
		return nil, apperrors.ErrProductInactive
	}

	if product.Stock < req.Quantity {
		return nil, apperrors.ErrInsufficientStock
	}

	cartItem, err := c.cartRepo.FindItemByCartAndProduct(cart.ID, req.ProductID)
	if err != nil {
		if errors.Is(err, apperrors.ErrCartItemNotFound) {
			// Item baru — create
			cartItem = &models.CartItem{
				CartID:    cart.ID,
				ProductID: req.ProductID,
				Quantity:  req.Quantity,
			}
			if err := c.cartRepo.CreateItem(cartItem); err != nil {
				return nil, apperrors.ErrInternal
			}
		} else {
			return nil, apperrors.ErrInternal
		}
	} else {
		newQuantity := cartItem.Quantity + req.Quantity
		if product.Stock < newQuantity {
			return nil, apperrors.ErrInsufficientStock
		}
		cartItem.Quantity = newQuantity
		if err := c.cartRepo.UpdateItem(cartItem); err != nil {
			return nil, apperrors.ErrInternal
		}
	}

	cart, err = c.cartRepo.FindByUserID(userID)
	if err != nil {
		return nil, apperrors.ErrInternal
	}

	response := dto.ToCartResponse(cart)
	return &response, nil
}

func (c *cartService) UpdateItem(userID uint, itemID uint, req dto.UpdateCartItemRequest) (*dto.CartResponse, error) {
	cart, err := c.cartRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrCartNotFound) {
			return nil, apperrors.ErrCartNotFound
		}
		return nil, apperrors.ErrInternal
	}

	cartItem, err := c.cartRepo.FindItemByID(itemID)
	if err != nil {
		if errors.Is(err, apperrors.ErrCartItemNotFound) {
			return nil, apperrors.ErrCartItemNotFound
		}
		return nil, apperrors.ErrInternal
	}

	if cartItem.CartID != cart.ID {
		return nil, apperrors.ErrForbidden
	}

	product, err := c.productRepo.FindByID(cartItem.ProductID)
	if err != nil {
		if errors.Is(err, apperrors.ErrProductNotFound) {
			return nil, apperrors.ErrProductNotFound
		}
		return nil, apperrors.ErrInternal
	}

	if !product.IsActive {
		return nil, apperrors.ErrProductInactive
	}

	if product.Stock < req.Quantity {
		return nil, apperrors.ErrInsufficientStock
	}

	cartItem.Quantity = req.Quantity
	if err := c.cartRepo.UpdateItem(cartItem); err != nil {
		return nil, apperrors.ErrInternal
	}

	cart, err = c.cartRepo.FindByUserID(userID)
	if err != nil {
		return nil, apperrors.ErrInternal
	}

	response := dto.ToCartResponse(cart)
	return &response, nil
}

func (c *cartService) RemoveItem(userID uint, itemID uint) (*dto.CartResponse, error) {
	cart, err := c.cartRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrCartNotFound) {
			return nil, apperrors.ErrCartNotFound
		}
		return nil, apperrors.ErrInternal
	}

	cartItem, err := c.cartRepo.FindItemByID(itemID)
	if err != nil {
		if errors.Is(err, apperrors.ErrCartItemNotFound) {
			return nil, apperrors.ErrCartItemNotFound
		}
		return nil, apperrors.ErrInternal
	}

	if cartItem.CartID != cart.ID {
		return nil, apperrors.ErrForbidden
	}

	if err := c.cartRepo.DeleteItem(itemID); err != nil {
		return nil, apperrors.ErrInternal
	}

	cart, err = c.cartRepo.FindByUserID(userID)
	if err != nil {
		return nil, apperrors.ErrInternal
	}

	response := dto.ToCartResponse(cart)
	return &response, nil
}
