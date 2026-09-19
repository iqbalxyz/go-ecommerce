package models

import "time"

type Order struct {
	ID         uint        `gorm:"primaryKey" json:"id"`
	UserID     uint        `gorm:"not null;index" json:"user_id"`
	TotalPrice int64       `gorm:"not null;default:0" json:"total_price"`
	Status     string      `gorm:"size:20;not null;default:'pending'" json:"status"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	Items      []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

type OrderItem struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	OrderID     uint      `gorm:"not null;index" json:"order_id"`
	ProductID   uint      `gorm:"not null" json:"product_id"`
	ProductName string    `gorm:"size:200;not null" json:"product_name"`
	Price       int64     `gorm:"not null" json:"price"`
	Quantity    int       `gorm:"not null" json:"quantity"`
	Subtotal    int64     `gorm:"not null" json:"subtotal"`
	CreatedAt   time.Time `json:"created_at"`
}

func (Order) TableName() string {
	return "orders"
}

func (OrderItem) TableName() string {
	return "order_items"
}
