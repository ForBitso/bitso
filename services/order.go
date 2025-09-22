package services

import (
	"errors"
	"fmt"
	"time"

	"go-shop/database"
	"go-shop/models"

	"gorm.io/gorm"
)

type OrderService struct{}

func NewOrderService() *OrderService {
	return &OrderService{}
}

func (os *OrderService) CreateOrder(userID uint, req *models.OrderCreateRequest) (*models.OrderResponse, error) {
	// Start transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Generate order number
	orderNumber := fmt.Sprintf("ORD-%d-%d", time.Now().Unix(), userID)

	// Calculate total amount and validate products
	var totalAmount float64
	var orderItems []models.OrderItem

	for _, item := range req.Items {
		// Get product details
		var product models.Product
		if err := tx.First(&product, item.ProductID).Error; err != nil {
			tx.Rollback()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("product not found")
			}
			return nil, errors.New("database error")
		}

		// Check stock
		if product.Stock < item.Quantity {
			tx.Rollback()
			return nil, fmt.Errorf("insufficient stock for product %s", product.Title)
		}

		// Calculate item total
		itemTotal := product.Price * float64(item.Quantity)
		totalAmount += itemTotal

		// Create order item
		orderItem := models.OrderItem{
			ProductID:     item.ProductID,
			Quantity:      item.Quantity,
			PriceAtMoment: product.Price,
		}
		orderItems = append(orderItems, orderItem)
	}

	// Create order
	order := models.Order{
		UserID:      userID,
		OrderNumber: orderNumber,
		Status:      models.OrderStatusPending,
		TotalAmount: totalAmount,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to create order")
	}

	// Create order items
	for i := range orderItems {
		orderItems[i].OrderID = order.ID
	}

	if err := tx.Create(&orderItems).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to create order items")
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("failed to commit transaction")
	}

	// Load order with items for response
	var orderWithItems models.Order
	if err := database.DB.Preload("OrderItems.Product").First(&orderWithItems, order.ID).Error; err != nil {
		return nil, errors.New("failed to load order")
	}

	// Build response
	var orderItemResponses []models.OrderItemResponse
	for _, item := range orderWithItems.OrderItems {
		orderItemResponses = append(orderItemResponses, models.OrderItemResponse{
			ID:            item.ID,
			OrderID:       item.OrderID,
			ProductID:     item.ProductID,
			Quantity:      item.Quantity,
			PriceAtMoment: item.PriceAtMoment,
		})
	}

	return &models.OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		OrderNumber: order.OrderNumber,
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
		OrderItems:  orderItemResponses,
	}, nil
}

func (os *OrderService) GetUserOrders(userID uint) ([]models.OrderResponse, error) {
	var orders []models.Order
	if err := database.DB.Preload("OrderItems").Where("user_id = ?", userID).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, errors.New("failed to get orders")
	}

	var orderResponses []models.OrderResponse
	for _, order := range orders {
		var orderItemResponses []models.OrderItemResponse
		for _, item := range order.OrderItems {
			orderItemResponses = append(orderItemResponses, models.OrderItemResponse{
				ID:            item.ID,
				OrderID:       item.OrderID,
				ProductID:     item.ProductID,
				Quantity:      item.Quantity,
				PriceAtMoment: item.PriceAtMoment,
			})
		}

		orderResponses = append(orderResponses, models.OrderResponse{
			ID:          order.ID,
			UserID:      order.UserID,
			OrderNumber: order.OrderNumber,
			Status:      order.Status,
			TotalAmount: order.TotalAmount,
			CreatedAt:   order.CreatedAt,
			UpdatedAt:   order.UpdatedAt,
			OrderItems:  orderItemResponses,
		})
	}

	return orderResponses, nil
}

func (os *OrderService) GetOrderByID(orderID, userID uint) (*models.OrderResponse, error) {
	var order models.Order
	if err := database.DB.Preload("OrderItems.Product").Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, errors.New("database error")
	}

	var orderItemResponses []models.OrderItemResponse
	for _, item := range order.OrderItems {
		orderItemResponses = append(orderItemResponses, models.OrderItemResponse{
			ID:            item.ID,
			OrderID:       item.OrderID,
			ProductID:     item.ProductID,
			Quantity:      item.Quantity,
			PriceAtMoment: item.PriceAtMoment,
		})
	}

	return &models.OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		OrderNumber: order.OrderNumber,
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
		OrderItems:  orderItemResponses,
	}, nil
}

func (os *OrderService) UpdateOrderStatus(orderID, userID uint, req *models.OrderUpdateRequest) (*models.OrderResponse, error) {
	var order models.Order
	if err := database.DB.Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, errors.New("database error")
	}

	// Users can only cancel their own orders
	if req.Status != models.OrderStatusCancelled {
		return nil, errors.New("users can only cancel orders")
	}

	// Check if order can be cancelled
	if order.Status == models.OrderStatusDelivered || order.Status == models.OrderStatusCancelled {
		return nil, errors.New("order cannot be cancelled")
	}

	// Update status to cancelled
	order.Status = req.Status

	if err := database.DB.Save(&order).Error; err != nil {
		return nil, errors.New("failed to update order")
	}

	return &models.OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		OrderNumber: order.OrderNumber,
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}, nil
}

// Admin functions
func (os *OrderService) GetAllOrders(limit, offset int) ([]models.OrderResponse, error) {
	var orders []models.Order
	if err := database.DB.Preload("OrderItems").Preload("User").Limit(limit).Offset(offset).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, errors.New("failed to get orders")
	}

	var orderResponses []models.OrderResponse
	for _, order := range orders {
		var orderItemResponses []models.OrderItemResponse
		for _, item := range order.OrderItems {
			orderItemResponses = append(orderItemResponses, models.OrderItemResponse{
				ID:            item.ID,
				OrderID:       item.OrderID,
				ProductID:     item.ProductID,
				Quantity:      item.Quantity,
				PriceAtMoment: item.PriceAtMoment,
			})
		}

		orderResponses = append(orderResponses, models.OrderResponse{
			ID:          order.ID,
			UserID:      order.UserID,
			OrderNumber: order.OrderNumber,
			Status:      order.Status,
			TotalAmount: order.TotalAmount,
			CreatedAt:   order.CreatedAt,
			UpdatedAt:   order.UpdatedAt,
			OrderItems:  orderItemResponses,
		})
	}

	return orderResponses, nil
}

func (os *OrderService) ConfirmOrder(orderID uint) (*models.OrderResponse, error) {
	// Start transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var order models.Order
	if err := tx.Preload("OrderItems").First(&order, orderID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, errors.New("database error")
	}

	// Check if order is paid
	if order.Status != models.OrderStatusPaid {
		tx.Rollback()
		return nil, errors.New("order must be paid before confirmation")
	}

	// Update stock and order_count for each item
	for _, item := range order.OrderItems {
		// Update stock
		if err := tx.Model(&models.Product{}).Where("id = ?", item.ProductID).Update("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("failed to update product stock")
		}

		// Update order_count (increment by 1 for each confirmed order)
		if err := tx.Model(&models.Product{}).Where("id = ?", item.ProductID).Update("order_count", gorm.Expr("order_count + 1")).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("failed to update product order count")
		}
	}

	// Update order status
	order.Status = models.OrderStatusConfirmed
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to update order status")
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("failed to commit transaction")
	}

	return &models.OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		OrderNumber: order.OrderNumber,
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}, nil
}

// ShipOrder marks an order as shipped (Admin/Seller only)
func (os *OrderService) ShipOrder(orderID uint) (*models.OrderResponse, error) {
	var order models.Order
	if err := database.DB.Where("id = ?", orderID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, errors.New("database error")
	}

	// Check if order can be shipped
	if order.Status != models.OrderStatusConfirmed {
		return nil, errors.New("order must be confirmed before shipping")
	}

	// Update status to shipped
	order.Status = models.OrderStatusShipped

	if err := database.DB.Save(&order).Error; err != nil {
		return nil, errors.New("failed to update order status")
	}

	return &models.OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		OrderNumber: order.OrderNumber,
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}, nil
}

// DeliverOrder marks an order as delivered (Admin only)
func (os *OrderService) DeliverOrder(orderID uint) (*models.OrderResponse, error) {
	var order models.Order
	if err := database.DB.Where("id = ?", orderID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, errors.New("database error")
	}

	// Check if order can be delivered
	if order.Status != models.OrderStatusShipped {
		return nil, errors.New("order must be shipped before delivery")
	}

	// Update status to delivered
	order.Status = models.OrderStatusDelivered

	if err := database.DB.Save(&order).Error; err != nil {
		return nil, errors.New("failed to update order status")
	}

	return &models.OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		OrderNumber: order.OrderNumber,
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}, nil
}

// CancelOrder cancels an order (User or Admin)
func (os *OrderService) CancelOrder(orderID uint) (*models.OrderResponse, error) {
	var order models.Order
	if err := database.DB.Where("id = ?", orderID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, errors.New("database error")
	}

	// Check if order can be cancelled
	if order.Status == models.OrderStatusDelivered || order.Status == models.OrderStatusCancelled {
		return nil, errors.New("order cannot be cancelled")
	}

	// Update status to cancelled
	order.Status = models.OrderStatusCancelled

	if err := database.DB.Save(&order).Error; err != nil {
		return nil, errors.New("failed to update order status")
	}

	return &models.OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		OrderNumber: order.OrderNumber,
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}, nil
}

// PayOrder marks an order as paid (User only)
func (os *OrderService) PayOrder(orderID, userID uint) (*models.OrderResponse, error) {
	var order models.Order
	if err := database.DB.Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, errors.New("database error")
	}

	// Check if order can be paid
	if order.Status != models.OrderStatusPending {
		return nil, errors.New("only pending orders can be paid")
	}

	// Update status to paid
	order.Status = models.OrderStatusPaid

	if err := database.DB.Save(&order).Error; err != nil {
		return nil, errors.New("failed to update order status")
	}

	return &models.OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		OrderNumber: order.OrderNumber,
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}, nil
}
