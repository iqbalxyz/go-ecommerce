package handler

import (
	"errors"
	apperrors "go-ecommerce/internal/errors"
	"go-ecommerce/internal/service"
	"go-ecommerce/internal/utils"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (oh *OrderHandler) Checkout(c fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}

	order, err := oh.orderService.Checkout(userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrEmptyCart) {
			return utils.Error(c, fiber.StatusBadRequest, "empty cart")
		}

		if errors.Is(err, apperrors.ErrInsufficientStock) {
			return utils.Error(c, fiber.StatusBadRequest, "insufficient stock")
		}

		if errors.Is(err, apperrors.ErrProductInactive) {
			return utils.Error(c, fiber.StatusBadRequest, "product inactive")
		}

		if errors.Is(err, apperrors.ErrInternal) {
			log.Printf("Order internal error: %v", err)
			return utils.Error(c, fiber.StatusInternalServerError, "internal server error")
		}

		log.Printf("Order checkout unexpected error: %v", err)
		return utils.Error(c, fiber.StatusInternalServerError, "internal server error")
	}

	return utils.Success(c, fiber.StatusOK, order)
}

func (oh *OrderHandler) List(c fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}

	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return utils.Error(c, fiber.StatusBadRequest, "invalid page")
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		return utils.Error(c, fiber.StatusBadRequest, "invalid limit")
	}

	offset := (page - 1) * limit

	res, count, err := oh.orderService.List(userID, limit, offset)
	if err != nil {
		if errors.Is(err, apperrors.ErrInternal) {
			log.Printf("Order internal error: %v", err)
			return utils.Error(c, fiber.StatusInternalServerError, "internal server error")
		}

		log.Printf("Order get unexpected error: %v", err)
		return utils.Error(c, fiber.StatusInternalServerError, "internal server error")
	}

	return utils.Success(c, fiber.StatusOK, fiber.Map{
		"data":       res,
		"total_data": count,
		"page":       page,
		"limit":      limit,
	})
}

func (oh *OrderHandler) GetByID(c fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}

	idStr := c.Params("id")
	orderID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || orderID < 1 {
		return utils.Error(c, fiber.StatusBadRequest, "invalid id")
	}

	order, err := oh.orderService.GetByID(userID, uint(orderID))
	if err != nil {
		if errors.Is(err, apperrors.ErrOrderNotFound) {
			return utils.Error(c, fiber.StatusNotFound, "order not found")
		}
		if errors.Is(err, apperrors.ErrInternal) {
			log.Printf("Order internal error: %v", err)
			return utils.Error(c, fiber.StatusInternalServerError, "internal server error")
		}

		log.Printf("Order get unexpected error: %v", err)
		return utils.Error(c, fiber.StatusInternalServerError, "internal server error")
	}

	return utils.Success(c, fiber.StatusOK, order)
}
