package services

import (
	"errors"

	"go-shop/database"
	"go-shop/models"

	"gorm.io/gorm"
)

type FavoriteService struct{}

func NewFavoriteService() *FavoriteService {
	return &FavoriteService{}
}

func (fs *FavoriteService) AddToFavorites(userID uint, req *models.FavoriteCreateRequest) (*models.FavoriteResponse, error) {
	// Check if already in favorites
	var existingFavorite models.Favorite
	if err := database.DB.Where("user_id = ? AND item_id = ? AND item_type = ?",
		userID, req.ItemID, req.ItemType).First(&existingFavorite).Error; err == nil {
		return nil, errors.New("item already in favorites")
	}

	favorite := models.Favorite{
		UserID:   userID,
		ItemID:   req.ItemID,
		ItemType: req.ItemType,
	}

	if err := database.DB.Create(&favorite).Error; err != nil {
		return nil, errors.New("failed to add to favorites")
	}

	return &models.FavoriteResponse{
		ID:        favorite.ID,
		UserID:    favorite.UserID,
		ItemID:    favorite.ItemID,
		ItemType:  favorite.ItemType,
		CreatedAt: favorite.CreatedAt,
	}, nil
}

func (fs *FavoriteService) GetUserFavorites(userID uint) ([]models.FavoriteResponse, error) {
	var favorites []models.Favorite
	if err := database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&favorites).Error; err != nil {
		return nil, errors.New("failed to get favorites")
	}

	var favoriteResponses []models.FavoriteResponse
	for _, favorite := range favorites {
		favoriteResponses = append(favoriteResponses, models.FavoriteResponse{
			ID:        favorite.ID,
			UserID:    favorite.UserID,
			ItemID:    favorite.ItemID,
			ItemType:  favorite.ItemType,
			CreatedAt: favorite.CreatedAt,
		})
	}

	return favoriteResponses, nil
}

func (fs *FavoriteService) RemoveFromFavorites(userID, favoriteID uint) error {
	var favorite models.Favorite
	if err := database.DB.Where("id = ? AND user_id = ?", favoriteID, userID).First(&favorite).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("favorite not found")
		}
		return errors.New("database error")
	}

	if err := database.DB.Delete(&favorite).Error; err != nil {
		return errors.New("failed to remove from favorites")
	}

	return nil
}

func (fs *FavoriteService) IsInFavorites(userID, itemID uint, itemType string) (bool, error) {
	var favorite models.Favorite
	if err := database.DB.Where("user_id = ? AND item_id = ? AND item_type = ?",
		userID, itemID, itemType).First(&favorite).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, errors.New("database error")
	}

	return true, nil
}
