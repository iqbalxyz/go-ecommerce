package service

import (
	"errors"
	"log"
	"time"

	"go-ecommerce/internal/dto"
	apperrors "go-ecommerce/internal/errors"
	"go-ecommerce/internal/models"
	"go-ecommerce/internal/repository"
	"go-ecommerce/internal/utils"
)

type UserService interface {
	Register(req dto.RegisterRequest) (*dto.UserResponse, error)
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
	GetByID(id uint) (*dto.UserResponse, error)
}

type userService struct {
	userRepo  repository.UserRepository
	jwtSecret string
	jwtExpiry time.Duration
}

func NewUserService(userRepo repository.UserRepository, jwtSecret string, jwtExpiry time.Duration) UserService {
	return &userService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
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

func (s *userService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Check if email is already registered
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return nil, apperrors.ErrInvalidCredentials
		}
		return nil, apperrors.ErrInternal
	}

	if !user.IsActive {
		return nil, apperrors.ErrInvalidCredentials
	}

	// Compare password
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, apperrors.ErrInvalidCredentials
	}

	// Generate JWT Token
	token, err := utils.GenerateToken(user.ID, user.Role, s.jwtSecret, s.jwtExpiry)
	if err != nil {
		log.Printf("Token generation error: %v", err)
		return nil, apperrors.ErrInternal
	}

	response := dto.ToLoginResponse(user, token)
	return &response, nil
}

func (s *userService) GetByID(id uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return nil, apperrors.ErrInvalidCredentials
		}
		return nil, apperrors.ErrInternal
	}
	response := dto.ToUserResponse(user)
	return &response, nil
}
