package models

import (
	"time"
)

type Cart struct {
	ID        uint
	UserID    uint `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Items     []CartItem `gorm:"foreignKey:CartID"`
}

type CartItem struct {
	ID        uint
	CartID    uint `gorm:"not null;index"`
	ProductID uint `gorm:"not null;index"`
	Quantity  int  `gorm:"not null;default:1"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Product   *Product `gorm:"foreignKey:ProductID"`
}

func (Cart) TableName() string {
	return "carts"
}

func (CartItem) TableName() string {
	return "cart_items"
}
