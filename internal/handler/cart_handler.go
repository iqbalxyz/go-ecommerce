package handler

import (
	"errors"
	"go-ecommerce/internal/dto"
	apperrors "go-ecommerce/internal/errors"
	"go-ecommerce/internal/service"
	"go-ecommerce/internal/utils"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type CartHandler struct {
	cartService service.CartService
}

func NewCartHandler(cartService service.CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

func (ch *CartHandler) GetCart(c fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}

	cart, err := ch.cartService.GetCart(userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrInternal) {
			log.Printf("Cart internal error: %v", err)
			return utils.Error(c, fiber.StatusInternalServerError, "internal server error")
		}
		log.Printf("Cart get unexpected error: %v", err)
		return utils.Error(c, fiber.StatusInternalServerError, "internal server error")
	}

	return utils.Success(c, fiber.StatusOK, cart)
}

func (ch *CartHandler) AddItem(c fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}

	var req dto.AddToCartRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "invalid request")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ValidationError(c, utils.FormatValidationErrors(err))
	}

	cart, err := ch.cartService.AddItem(userID, req)
	if err != nil {
		if errors.Is(err, apperrors.ErrProductNotFound) {
			return utils.Error(c, fiber.StatusNotFound, "product not found")
		}

		if errors.Is(err, apperrors.ErrCartNotFound) {
			return utils.Error(c, fiber.StatusNotFound, "cart not found")
		}

		if errors.Is(err, apperrors.ErrInternal) {
			log.Printf("Cart internal error: %v", err)
			return utils.Error(c, fiber.StatusInternalServerError, "internal server error")
		}

		if errors.Is(err, apperrors.ErrInsufficientStock) {
			return utils.Error(c, fiber.StatusBadRequest, "insufficient stock")
		}

		if errors.Is(err, apperrors.ErrProductInactive) {
			return utils.Error(c, fiber.StatusBadRequest, "product inactive")
		}

		log.Printf("Cart add to cart unexpected error: %v", err)
		return utils.Error(c, fiber.StatusInternalServerError, "internal server error")
	}

	return utils.Success(c, fiber.StatusCreated, cart)
}

func (ch *CartHandler) UpdateItem(c fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}

	id := c.Params("id")
	parsedID, err := strconv.ParseUint(id, 10, 64)
	if err != nil || parsedID < 1 {
		return utils.Error(c, fiber.StatusBadRequest, "invalid id")
	}

	var req dto.UpdateCartItemRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "invalid request")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ValidationError(c, utils.FormatValidationErrors(err))
	}

	cart, err := ch.cartService.UpdateItem(userID, uint(parsedID), req)
	if err != nil {
		if errors.Is(err, apperrors.ErrCartNotFound) {
			return utils.Error(c, fiber.StatusNotFound, "cart not found")
		}
		if errors.Is(err, apperrors.ErrCartItemNotFound) {
			return utils.Error(c, fiber.StatusNotFound, "cart item not found")
		}
		if errors.Is(err, apperrors.ErrForbidden) {
			return utils.Error(c, fiber.StatusForbidden, "forbidden")
		}
		if errors.Is(err, apperrors.ErrProductNotFound) {
			return utils.Error(c, fiber.StatusNotFound, "product not found")
		}
		if errors.Is(err, apperrors.ErrProductInactive) {
			return utils.Error(c, fiber.StatusBadRequest, "product inactive")
		}
		if errors.Is(err, apperrors.ErrInsufficientStock) {
			return utils.Error(c, fiber.StatusBadRequest, "insufficient stock")
		}

		if errors.Is(err, apperrors.ErrInternal) {
			log.Printf("Cart internal error: %v", err)
			return utils.Error(c, fiber.StatusInternalServerError, "internal server error")
		}
		log.Printf("Cart update unexpected error: %v", err)
		return utils.Error(c, fiber.StatusInternalServerError, "internal server error")
	}

	return utils.Success(c, fiber.StatusOK, cart)
}

func (ch *CartHandler) RemoveItem(c fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}

	id := c.Params("id")
	parsedID, err := strconv.ParseUint(id, 10, 64)
	if err != nil || parsedID < 1 {
		return utils.Error(c, fiber.StatusBadRequest, "invalid id")
	}

	_, err = ch.cartService.RemoveItem(userID, uint(parsedID))
	if err != nil {
		if errors.Is(err, apperrors.ErrCartNotFound) {
			return utils.Error(c, fiber.StatusNotFound, "cart not found")
		}
		if errors.Is(err, apperrors.ErrCartItemNotFound) {
			return utils.Error(c, fiber.StatusNotFound, "cart item not found")
		}
		if errors.Is(err, apperrors.ErrForbidden) {
			return utils.Error(c, fiber.StatusForbidden, "forbidden")
		}
		if errors.Is(err, apperrors.ErrProductNotFound) {
			return utils.Error(c, fiber.StatusNotFound, "product not found")
		}
		if errors.Is(err, apperrors.ErrProductInactive) {
			return utils.Error(c, fiber.StatusBadRequest, "product inactive")
		}
		if errors.Is(err, apperrors.ErrInsufficientStock) {
			return utils.Error(c, fiber.StatusBadRequest, "insufficient stock")
		}

		if errors.Is(err, apperrors.ErrInternal) {
			log.Printf("Cart internal error: %v", err)
			return utils.Error(c, fiber.StatusInternalServerError, "internal server error")
		}
		log.Printf("Cart update unexpected error: %v", err)
		return utils.Error(c, fiber.StatusInternalServerError, "internal server error")
	}

	return utils.Success(c, fiber.StatusOK, "item removed successfully")

}
