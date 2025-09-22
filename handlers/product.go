package handlers

import (
	"net/http"
	"strconv"

	"go-shop/models"
	"go-shop/services"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService *services.ProductService
}

func NewProductHandler(productService *services.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// GetProducts godoc
// @Summary Get products
// @Description Get products with optional filtering
// @Tags products
// @Accept json
// @Produce json
// @Param category_id query int false "Filter by category ID"
// @Param limit query int false "Limit results" default(20)
// @Param offset query int false "Offset results" default(0)
// @Success 200 {object} models.SuccessResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /products [get]
func (ph *ProductHandler) GetProducts(c *gin.Context) {
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

	products, err := ph.productService.GetProducts(categoryID, limit, offset)
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

// GetProductByID godoc
// @Summary Get product by ID
// @Description Get specific product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /products/{id} [get]
func (ph *ProductHandler) GetProductByID(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid product ID",
			Message: err.Error(),
		})
		return
	}

	product, err := ph.productService.GetProductByID(uint(productID))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "Product not found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Product retrieved successfully",
		Data:    product,
	})
}

// SearchProducts godoc
// @Summary Search products
// @Description Search products with filters and sorting
// @Tags products
// @Accept json
// @Produce json
// @Param title query string false "Search by title"
// @Param category_id query int false "Filter by category ID"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param sort_by query string false "Sort by: price_asc, price_desc, popularity_asc, popularity_desc, created_at_asc, created_at_desc"
// @Param limit query int false "Limit results" default(20)
// @Param offset query int false "Offset results" default(0)
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /products/search [get]
func (ph *ProductHandler) SearchProducts(c *gin.Context) {
	var req models.ProductSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request parameters",
			Message: err.Error(),
		})
		return
	}

	// Get user ID for logging (optional)
	var userID *uint
	if uid, exists := c.Get("user_id"); exists {
		userIDUint := uid.(uint)
		userID = &userIDUint
	}

	// Search products
	products, total, err := ph.productService.SearchProducts(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to search products",
			Message: err.Error(),
		})
		return
	}

	// Log search query
	if userID != nil || req.Title != "" {
		filters := models.JSONB{
			"category_id": req.CategoryID,
			"min_price":   req.MinPrice,
			"max_price":   req.MaxPrice,
			"sort_by":     req.SortBy,
		}
		ph.productService.LogSearch(userID, req.Title, filters, len(products))
	}

	// Prepare response
	response := gin.H{
		"products": products,
		"total":    total,
		"limit":    req.Limit,
		"offset":   req.Offset,
		"has_more": int64(req.Offset+req.Limit) < total,
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Products found successfully",
		Data:    response,
	})
}
