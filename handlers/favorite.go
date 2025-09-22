package handlers

import (
	"net/http"
	"strconv"

	"go-shop/models"
	"go-shop/services"

	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	favoriteService *services.FavoriteService
}

func NewFavoriteHandler(favoriteService *services.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{
		favoriteService: favoriteService,
	}
}

// AddToFavorites godoc
// @Summary Add item to favorites
// @Description Add an item to user's favorites
// @Tags favorites
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.FavoriteCreateRequest true "Favorite item data"
// @Success 201 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /favorites [post]
func (fh *FavoriteHandler) AddToFavorites(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	var req models.FavoriteCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request data",
			Message: err.Error(),
		})
		return
	}

	favorite, err := fh.favoriteService.AddToFavorites(userID.(uint), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to add to favorites",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Message: "Item added to favorites successfully",
		Data:    favorite,
	})
}

// GetUserFavorites godoc
// @Summary Get user favorites
// @Description Get all favorite items for the authenticated user
// @Tags favorites
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.SuccessResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /favorites [get]
func (fh *FavoriteHandler) GetUserFavorites(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	favorites, err := fh.favoriteService.GetUserFavorites(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get favorites",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Favorites retrieved successfully",
		Data:    favorites,
	})
}

// RemoveFromFavorites godoc
// @Summary Remove item from favorites
// @Description Remove an item from user's favorites
// @Tags favorites
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Favorite ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /favorites/{id} [delete]
func (fh *FavoriteHandler) RemoveFromFavorites(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	favoriteIDStr := c.Param("id")
	favoriteID, err := strconv.ParseUint(favoriteIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid favorite ID",
			Message: err.Error(),
		})
		return
	}

	err = fh.favoriteService.RemoveFromFavorites(userID.(uint), uint(favoriteID))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "Failed to remove from favorites",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Item removed from favorites successfully",
	})
}

// CheckFavorite godoc
// @Summary Check if item is in favorites
// @Description Check if a specific item is in user's favorites
// @Tags favorites
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param item_id query int true "Item ID"
// @Param item_type query string true "Item Type"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /favorites/check [get]
func (fh *FavoriteHandler) CheckFavorite(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	itemIDStr := c.Query("item_id")
	itemID, err := strconv.ParseUint(itemIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid item ID",
			Message: err.Error(),
		})
		return
	}

	itemType := c.Query("item_type")
	if itemType == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Item type is required",
		})
		return
	}

	isFavorite, err := fh.favoriteService.IsInFavorites(userID.(uint), uint(itemID), itemType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to check favorite status",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Favorite status checked successfully",
		Data: gin.H{
			"is_favorite": isFavorite,
		},
	})
}
