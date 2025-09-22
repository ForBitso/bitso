package models

import (
	"time"

	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"   // Создан, но не оплачен
	OrderStatusPaid      OrderStatus = "paid"      // Оплачен, но не подтвержден
	OrderStatusConfirmed OrderStatus = "confirmed" // Админ подтвердил
	OrderStatusShipped   OrderStatus = "shipped"   // Отправлен
	OrderStatusDelivered OrderStatus = "delivered" // Доставлен
	OrderStatusCancelled OrderStatus = "cancelled" // Отменен
)

type Order struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null"`
	OrderNumber string         `json:"order_number" gorm:"uniqueIndex;not null"`
	Status      OrderStatus    `json:"status" gorm:"default:'pending'"`
	TotalAmount float64        `json:"total_amount" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	User       User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
	OrderItems []OrderItem `json:"order_items,omitempty" gorm:"foreignKey:OrderID"`
}

type OrderCreateRequest struct {
	Items []OrderItemRequest `json:"items" binding:"required,min=1"`
}

type OrderUpdateRequest struct {
	Status OrderStatus `json:"status" binding:"required,oneof=pending paid confirmed shipped delivered cancelled"`
}

type OrderResponse struct {
	ID          uint                `json:"id"`
	UserID      uint                `json:"user_id"`
	OrderNumber string              `json:"order_number"`
	Status      OrderStatus         `json:"status"`
	TotalAmount float64             `json:"total_amount"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	OrderItems  []OrderItemResponse `json:"order_items,omitempty"`
}

type OrderItem struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	OrderID       uint           `json:"order_id" gorm:"not null"`
	ProductID     uint           `json:"product_id" gorm:"not null"`
	Quantity      int            `json:"quantity" gorm:"not null"`
	PriceAtMoment float64        `json:"price_at_moment" gorm:"not null"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Order   Order   `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	Product Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

type OrderItemRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

type OrderItemResponse struct {
	ID            uint             `json:"id"`
	OrderID       uint             `json:"order_id"`
	ProductID     uint             `json:"product_id"`
	Quantity      int              `json:"quantity"`
	PriceAtMoment float64          `json:"price_at_moment"`
	Product       *ProductResponse `json:"product,omitempty"`
}
