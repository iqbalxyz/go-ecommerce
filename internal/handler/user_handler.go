package handler

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v3"

	"go-ecommerce/internal/dto"
	apperrors "go-ecommerce/internal/errors"
	"go-ecommerce/internal/service"
	"go-ecommerce/internal/utils"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "invalid request body")

	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ValidationError(c, utils.FormatValidationErrors(err))
	}

	response, err := h.userService.Register(req)
	if err != nil {
		if errors.Is(err, apperrors.ErrEmailExists) {
			log.Printf("Register email exists error: %v", err)
			return utils.Error(c, fiber.StatusConflict, "email already exists")
		}
		if errors.Is(err, apperrors.ErrInternal) {
			log.Printf("Register internal server error: %v", err)
			return utils.Error(c, fiber.StatusInternalServerError, "internal server error")
		}
		log.Printf("Register unexpected server error: %v", err)
		return utils.Error(c, fiber.StatusInternalServerError, "internal server error")
	}

	return utils.Success(c, fiber.StatusCreated, response)
}
