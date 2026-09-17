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

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (ph *ProductHandler) Create(c fiber.Ctx) error {
	var req dto.CreateProductRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "invalid request")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ValidationError(c, utils.FormatValidationErrors(err))
	}

	res, err := ph.productService.Create(req)
	if err != nil {
		if errors.Is(err, apperrors.ErrProductExists) {
			return utils.Error(c, fiber.StatusConflict, "product already exists")
		}

		if errors.Is(err, apperrors.ErrInternal) {
			log.Printf("Product internal error: %v", err)
			return utils.Error(c, fiber.StatusInternalServerError, "internal server error")
		}

		log.Printf("Product create unexpected error: %v", err)
		return utils.Error(c, fiber.StatusInternalServerError, "internal server error")
	}

	return utils.Success(c, fiber.StatusCreated, res)
}

func (ph *ProductHandler) List(c fiber.Ctx) error {
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

	res, count, err := ph.productService.List(limit, offset)
	if err != nil {
		if errors.Is(err, apperrors.ErrInternal) {
			log.Printf("Product internal error: %v", err)
			return utils.Error(c, fiber.StatusInternalServerError, "internal server error")
		}

		log.Printf("Product list unexpected error: %v", err)
		return utils.Error(c, fiber.StatusInternalServerError, "internal server error")
	}

	return utils.Success(c, fiber.StatusOK, fiber.Map{
		"data": res,
		"meta": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": count,
		},
	})
}

func (ph *ProductHandler) GetByID(c fiber.Ctx) error {
	id := c.Params("id")
	parsedID, err := strconv.ParseUint(id, 10, 64)
	if err != nil || parsedID < 1 {
		return utils.Error(c, fiber.StatusBadRequest, "invalid id")
	}

	res, err := ph.productService.GetByID(uint(parsedID))
	if err != nil {
		if errors.Is(err, apperrors.ErrProductNotFound) {
			return utils.Error(c, fiber.StatusNotFound, "product not found")
		}

		if errors.Is(err, apperrors.ErrInternal) {
			log.Printf("Product internal error: %v", err)
			return utils.Error(c, fiber.StatusInternalServerError, "internal server error")
		}

		log.Printf("Product get by id unexpected error: %v", err)
		return utils.Error(c, fiber.StatusInternalServerError, "internal server error")
	}

	return utils.Success(c, fiber.StatusOK, res)
}

func (ph *ProductHandler) Update(c fiber.Ctx) error {
	id := c.Params("id")
	parsedID, err := strconv.ParseUint(id, 10, 64)
	if err != nil || parsedID < 1 {
		return utils.Error(c, fiber.StatusBadRequest, "invalid id")
	}

	var req dto.CreateProductRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "invalid request")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ValidationError(c, utils.FormatValidationErrors(err))
	}
	res, err := ph.productService.Update(uint(parsedID), req)
	if err != nil {
		if errors.Is(err, apperrors.ErrProductNotFound) {
			return utils.Error(c, fiber.StatusNotFound, "product not found")
		}

		if errors.Is(err, apperrors.ErrInternal) {
			log.Printf("Product internal error: %v", err)
			return utils.Error(c, fiber.StatusInternalServerError, "internal server error")
		}

		log.Printf("Product update unexpected error: %v", err)
		return utils.Error(c, fiber.StatusInternalServerError, "internal server error")
	}

	return utils.Success(c, fiber.StatusOK, res)
}

func (ph *ProductHandler) Delete(c fiber.Ctx) error {
	id := c.Params("id")
	parsedID, err := strconv.ParseUint(id, 10, 64)
	if err != nil || parsedID < 1 {
		return utils.Error(c, fiber.StatusBadRequest, "invalid id")
	}

	err = ph.productService.Delete(uint(parsedID))
	if err != nil {
		if errors.Is(err, apperrors.ErrProductNotFound) {
			return utils.Error(c, fiber.StatusNotFound, "product not found")
		}

		if errors.Is(err, apperrors.ErrInternal) {
			log.Printf("Product internal error: %v", err)
			return utils.Error(c, fiber.StatusInternalServerError, "internal server error")
		}

		log.Printf("Product delete unexpected error: %v", err)
		return utils.Error(c, fiber.StatusInternalServerError, "internal server error")
	}

	return utils.Success(c, fiber.StatusOK, "product deleted successfully")
}
