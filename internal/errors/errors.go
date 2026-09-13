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
	ErrInvalidInput       = errors.New("invalid input")
	ErrInternal           = errors.New("internal server error")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrInvalidRequest     = errors.New("invalid request")
)
