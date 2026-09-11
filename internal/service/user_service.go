package service

import (
	"errors"
	"log"

	"go-ecommerce/internal/dto"
	apperrors "go-ecommerce/internal/errors"
	"go-ecommerce/internal/models"
	"go-ecommerce/internal/repository"
	"go-ecommerce/internal/utils"
)

type UserService interface {
	Register(req dto.RegisterRequest) (*dto.UserResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Register(req dto.RegisterRequest) (*dto.UserResponse, error) {
	// Check if email is already registered
	_, err := s.userRepo.FindByEmail(req.Email)
	if err == nil {
		return nil, apperrors.ErrEmailExists
	}
	if !errors.Is(err, apperrors.ErrUserNotFound) {
		return nil, apperrors.ErrInternal
	}

	// Hash password
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Printf("Register hash error: %v", err)
		return nil, apperrors.ErrInternal
	}

	// Create model
	user := models.User{
		Email:        req.Email,
		PasswordHash: hash,
		Name:         req.Name,
		Role:         "customer",
		IsActive:     true,
	}

	// Save user to database
	if err := s.userRepo.Create(&user); err != nil {
		log.Printf("Register create error: %v", err)
		return nil, apperrors.ErrInternal
	}

	// Return DTO response
	response := dto.ToUserResponse(&user)
	return &response, nil
}
