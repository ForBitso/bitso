package handlers

import (
	"net/http"
	"strconv"

	"go-shop/models"
	"go-shop/services"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService *services.CategoryService
}

func NewCategoryHandler(categoryService *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// GetCategories godoc
// @Summary Get categories
// @Description Get all product categories
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /categories [get]
func (ch *CategoryHandler) GetCategories(c *gin.Context) {
	categories, err := ch.categoryService.GetCategories()
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

// GetCategoryByID godoc
// @Summary Get category by ID
// @Description Get specific category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /categories/{id} [get]
func (ch *CategoryHandler) GetCategoryByID(c *gin.Context) {
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid category ID",
			Message: err.Error(),
		})
		return
	}

	category, err := ch.categoryService.GetCategoryByID(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "Category not found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Category retrieved successfully",
		Data:    category,
	})
}
