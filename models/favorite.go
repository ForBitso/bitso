package models

import (
	"time"

	"gorm.io/gorm"
)

type Favorite struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	ItemID    uint           `json:"item_id" gorm:"not null"`
	ItemType  string         `json:"item_type" gorm:"not null"` // product, category, etc.
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type FavoriteCreateRequest struct {
	ItemID   uint   `json:"item_id" binding:"required"`
	ItemType string `json:"item_type" binding:"required,oneof=product category"`
}

type FavoriteResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	ItemID    uint      `json:"item_id"`
	ItemType  string    `json:"item_type"`
	CreatedAt time.Time `json:"created_at"`
}
