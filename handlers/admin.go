package handlers

import (
	"net/http"
	"strconv"

	"go-shop/middleware"
	"go-shop/models"
	"go-shop/services"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	categoryService *services.CategoryService
	productService  *services.ProductService
	orderService    *services.OrderService
}

func NewAdminHandler(categoryService *services.CategoryService, productService *services.ProductService, orderService *services.OrderService) *AdminHandler {
	return &AdminHandler{
		categoryService: categoryService,
		productService:  productService,
		orderService:    orderService,
	}
}

// Category Management

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new product category (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CategoryCreateRequest true "Category creation data"
// @Success 201 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /admin/categories [post]
func (ah *AdminHandler) CreateCategory(c *gin.Context) {
	var req models.CategoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request data",
			Message: err.Error(),
		})
		return
	}

	category, err := ah.categoryService.CreateCategory(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to create category",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Message: "Category created successfully",
		Data:    category,
	})
}

// GetCategories godoc
// @Summary Get all categories
// @Description Get all product categories (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.SuccessResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /admin/categories [get]
func (ah *AdminHandler) GetCategories(c *gin.Context) {
	categories, err := ah.categoryService.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get categories",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Categories retrieved successfully",
		Data:    categories,
	})
}

// UpdateCategory godoc
// @Summary Update category
// @Description Update a product category (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Param request body models.CategoryUpdateRequest true "Category update data"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /admin/categories/{id} [put]
func (ah *AdminHandler) UpdateCategory(c *gin.Context) {
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid category ID",
			Message: err.Error(),
		})
		return
	}

	var req models.CategoryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request data",
			Message: err.Error(),
		})
		return
	}

	category, err := ah.categoryService.UpdateCategory(uint(categoryID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to update category",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Category updated successfully",
		Data:    category,
	})
}

// DeleteCategory godoc
// @Summary Delete category
// @Description Delete a product category (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /admin/categories/{id} [delete]
func (ah *AdminHandler) DeleteCategory(c *gin.Context) {
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid category ID",
			Message: err.Error(),
		})
		return
	}

	err = ah.categoryService.DeleteCategory(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to delete category",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Category deleted successfully",
	})
}

// Product Management

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.ProductCreateRequest true "Product creation data"
// @Success 201 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /admin/products [post]
func (ah *AdminHandler) CreateProduct(c *gin.Context) {
	var req models.ProductCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request data",
			Message: err.Error(),
		})
		return
	}

	product, err := ah.productService.CreateProduct(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to create product",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Message: "Product created successfully",
		Data:    product,
	})
}

// GetProducts godoc
// @Summary Get all products
// @Description Get all products with optional filtering (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category_id query int false "Filter by category ID"
// @Param limit query int false "Limit results" default(20)
// @Param offset query int false "Offset results" default(0)
// @Success 200 {object} models.SuccessResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /admin/products [get]
func (ah *AdminHandler) GetProducts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")
	categoryIDStr := c.Query("category_id")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	var categoryID *uint
	if categoryIDStr != "" {
		if id, err := strconv.ParseUint(categoryIDStr, 10, 32); err == nil {
			categoryIDUint := uint(id)
			categoryID = &categoryIDUint
		}
	}

	products, err := ah.productService.GetProducts(categoryID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get products",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Products retrieved successfully",
		Data:    products,
	})
}

// UpdateProduct godoc
// @Summary Update product
// @Description Update a product (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param request body models.ProductUpdateRequest true "Product update data"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /admin/products/{id} [put]
func (ah *AdminHandler) UpdateProduct(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid product ID",
			Message: err.Error(),
		})
		return
	}

	var req models.ProductUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request data",
			Message: err.Error(),
		})
		return
	}

	product, err := ah.productService.UpdateProduct(uint(productID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to update product",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Product updated successfully",
		Data:    product,
	})
}

// DeleteProduct godoc
// @Summary Delete product
// @Description Delete a product (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /admin/products/{id} [delete]
func (ah *AdminHandler) DeleteProduct(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid product ID",
			Message: err.Error(),
		})
		return
	}

	err = ah.productService.DeleteProduct(uint(productID))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to delete product",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Product deleted successfully",
	})
}

// Order Management

// GetAllOrders godoc
// @Summary Get all orders
// @Description Get all orders (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit results" default(20)
// @Param offset query int false "Offset results" default(0)
// @Success 200 {object} models.SuccessResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /admin/orders [get]
func (ah *AdminHandler) GetAllOrders(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	orders, err := ah.orderService.GetAllOrders(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get orders",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Orders retrieved successfully",
		Data:    orders,
	})
}

// ConfirmOrder godoc
// @Summary Confirm order
// @Description Confirm a paid order and update stock (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /admin/orders/{id}/confirm [post]
func (ah *AdminHandler) ConfirmOrder(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid order ID",
			Message: err.Error(),
		})
		return
	}

	// Get current user ID for logging
	currentUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	order, err := ah.orderService.ConfirmOrder(uint(orderID))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to confirm order",
			Message: err.Error(),
		})
		return
	}

	// Log sensitive operation
	middleware.LogSensitiveOperation("ORDER_CONFIRMED", currentUserID.(uint),
		"Order ID: "+orderIDStr+", Total Amount: "+strconv.FormatFloat(order.TotalAmount, 'f', 2, 64))

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Order confirmed successfully",
		Data:    order,
	})
}

// ShipOrder godoc
// @Summary Ship order
// @Description Mark order as shipped (Admin/Seller only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /admin/orders/{id}/ship [post]
func (ah *AdminHandler) ShipOrder(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid order ID",
			Message: err.Error(),
		})
		return
	}

	// Get current user ID for logging
	currentUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	order, err := ah.orderService.ShipOrder(uint(orderID))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to ship order",
			Message: err.Error(),
		})
		return
	}

	// Log sensitive operation
	middleware.LogSensitiveOperation("ORDER_SHIPPED", currentUserID.(uint),
		"Order ID: "+orderIDStr+", Total Amount: "+strconv.FormatFloat(order.TotalAmount, 'f', 2, 64))

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Order shipped successfully",
		Data:    order,
	})
}

// DeliverOrder godoc
// @Summary Deliver order
// @Description Mark order as delivered (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /admin/orders/{id}/deliver [post]
func (ah *AdminHandler) DeliverOrder(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid order ID",
			Message: err.Error(),
		})
		return
	}

	// Get current user ID for logging
	currentUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	order, err := ah.orderService.DeliverOrder(uint(orderID))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to deliver order",
			Message: err.Error(),
		})
		return
	}

	// Log sensitive operation
	middleware.LogSensitiveOperation("ORDER_DELIVERED", currentUserID.(uint),
		"Order ID: "+orderIDStr+", Total Amount: "+strconv.FormatFloat(order.TotalAmount, 'f', 2, 64))

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Order delivered successfully",
		Data:    order,
	})
}

// CancelOrder godoc
// @Summary Cancel order
// @Description Cancel an order (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /admin/orders/{id}/cancel [post]
func (ah *AdminHandler) CancelOrder(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid order ID",
			Message: err.Error(),
		})
		return
	}

	// Get current user ID for logging
	currentUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	order, err := ah.orderService.CancelOrder(uint(orderID))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to cancel order",
			Message: err.Error(),
		})
		return
	}

	// Log sensitive operation
	middleware.LogSensitiveOperation("ORDER_CANCELLED", currentUserID.(uint),
		"Order ID: "+orderIDStr+", Total Amount: "+strconv.FormatFloat(order.TotalAmount, 'f', 2, 64))

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Order cancelled successfully",
		Data:    order,
	})
}
