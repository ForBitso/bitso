package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// JSONB type for storing JSON data in PostgreSQL
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte(value.(string)), j)
	}
	return json.Unmarshal(bytes, j)
}

// StringArray type for storing array of strings
type StringArray []string

// Value implements the driver.Valuer interface
func (s StringArray) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

// Scan implements the sql.Scanner interface
func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte(value.(string)), s)
	}
	return json.Unmarshal(bytes, s)
}

type Product struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	CategoryID  *uint          `json:"category_id" gorm:"index"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description"`
	Images      StringArray    `json:"images" gorm:"type:jsonb"`
	Price       float64        `json:"price" gorm:"not null"`
	Model       string         `json:"model"`
	ExtraInfo   JSONB          `json:"extra_info" gorm:"type:jsonb"`
	Stock       int            `json:"stock" gorm:"not null;default:0"`
	OrderCount  int            `json:"order_count" gorm:"not null;default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Category   *Category   `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	OrderItems []OrderItem `json:"order_items,omitempty" gorm:"foreignKey:ProductID"`
}

type ProductCreateRequest struct {
	CategoryID  uint     `json:"category_id" binding:"required"`
	Title       string   `json:"title" binding:"required,min=2,max=200"`
	Description string   `json:"description" binding:"max=1000"`
	Images      []string `json:"images"`
	Price       float64  `json:"price" binding:"required,min=0"`
	Model       string   `json:"model" binding:"max=100"`
	ExtraInfo   JSONB    `json:"extra_info"`
	Stock       int      `json:"stock" binding:"min=0"`
}

type ProductUpdateRequest struct {
	CategoryID  *uint    `json:"category_id" binding:"omitempty"`
	Title       string   `json:"title" binding:"omitempty,min=2,max=200"`
	Description string   `json:"description" binding:"omitempty,max=1000"`
	Images      []string `json:"images"`
	Price       *float64 `json:"price" binding:"omitempty,min=0"`
	Model       string   `json:"model" binding:"omitempty,max=100"`
	ExtraInfo   JSONB    `json:"extra_info"`
	Stock       *int     `json:"stock" binding:"omitempty,min=0"`
}

type ProductResponse struct {
	ID          uint              `json:"id"`
	CategoryID  *uint             `json:"category_id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Images      []string          `json:"images"`
	Price       float64           `json:"price"`
	Model       string            `json:"model"`
	ExtraInfo   JSONB             `json:"extra_info"`
	Stock       int               `json:"stock"`
	OrderCount  int               `json:"order_count"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Category    *CategoryResponse `json:"category,omitempty"`
}

// Search request models
type ProductSearchRequest struct {
	Title      string   `form:"title"`
	CategoryID *uint    `form:"category_id"`
	MinPrice   *float64 `form:"min_price"`
	MaxPrice   *float64 `form:"max_price"`
	SortBy     string   `form:"sort_by" binding:"omitempty,oneof=price_asc price_desc popularity_asc popularity_desc created_at_asc created_at_desc"`
	Limit      int      `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset     int      `form:"offset" binding:"omitempty,min=0"`
}

// Search log model
type SearchLog struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    *uint     `json:"user_id" gorm:"index"`
	Query     string    `json:"query" gorm:"not null"`
	Filters   JSONB     `json:"filters" gorm:"type:jsonb"`
	Results   int       `json:"results"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
