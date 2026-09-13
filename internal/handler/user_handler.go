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

func (h *UserHandler) Login(c fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ValidationError(c, utils.FormatValidationErrors(err))
	}

	response, err := h.userService.Login(req)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidCredentials) {
			return utils.Error(c, fiber.StatusUnauthorized, "invalid credentials")
		}
		if errors.Is(err, apperrors.ErrInternal) {
			log.Printf("Login internal error: %v", err)
			return utils.Error(c, fiber.StatusInternalServerError, "internal server error")
		}
		log.Printf("Login unexpected error: %v", err)
		return utils.Error(c, fiber.StatusInternalServerError, "internal server error")
	}

	return utils.Success(c, fiber.StatusOK, response)
}

func (h *UserHandler) Me(c fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}

	response, err := h.userService.GetByID(userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return utils.Error(c, fiber.StatusNotFound, "user not found")
		}
		if errors.Is(err, apperrors.ErrInternal) {
			log.Printf("Me internal error: %v", err)
			return utils.Error(c, fiber.StatusInternalServerError, "internal server error")
		}
		log.Printf("Me unexpected error: %v", err)
		return utils.Error(c, fiber.StatusInternalServerError, "internal server error")
	}

	return utils.Success(c, fiber.StatusOK, response)
}
