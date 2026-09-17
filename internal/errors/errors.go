package errors

import "errors"

/*
Tujuan:
Definisikan error domain — error yang bermakna di level bisnis.
Ini yang akan dipakai service untuk bilang "email sudah ada" atau "user tidak ditemukan",
tanpa mencampur dengan detail teknis GORM/SQL.
*/

var (
	ErrEmailExists        = errors.New("email already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInternal           = errors.New("internal server error")
	ErrProductNotFound    = errors.New("product not found")
	ErrForbidden          = errors.New("forbidden")
	ErrProductExists      = errors.New("product already exists")
	ErrCartNotFound       = errors.New("cart not found")
	ErrCartItemNotFound   = errors.New("item not found")
	ErrProductInactive    = errors.New("product is inactive")
	ErrInsufficientStock  = errors.New("insufficient stock")
)
