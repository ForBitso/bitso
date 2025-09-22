package services

import (
	"errors"

	"go-shop/database"
	"go-shop/models"

	"gorm.io/gorm"
)

type CategoryService struct{}

func NewCategoryService() *CategoryService {
	return &CategoryService{}
}

func (cs *CategoryService) CreateCategory(req *models.CategoryCreateRequest) (*models.CategoryResponse, error) {
	category := models.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := database.DB.Create(&category).Error; err != nil {
		return nil, errors.New("failed to create category")
	}

	return &models.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil
}

func (cs *CategoryService) GetCategories() ([]models.CategoryResponse, error) {
	var categories []models.Category
	if err := database.DB.Find(&categories).Error; err != nil {
		return nil, errors.New("failed to get categories")
	}

	var categoryResponses []models.CategoryResponse
	for _, category := range categories {
		categoryResponses = append(categoryResponses, models.CategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			Description: category.Description,
			CreatedAt:   category.CreatedAt,
			UpdatedAt:   category.UpdatedAt,
		})
	}

	return categoryResponses, nil
}

func (cs *CategoryService) GetCategoryByID(categoryID uint) (*models.CategoryResponse, error) {
	var category models.Category
	if err := database.DB.First(&category, categoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, errors.New("database error")
	}

	return &models.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil
}

func (cs *CategoryService) UpdateCategory(categoryID uint, req *models.CategoryUpdateRequest) (*models.CategoryResponse, error) {
	var category models.Category
	if err := database.DB.First(&category, categoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, errors.New("database error")
	}

	// Update fields
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}

	if err := database.DB.Save(&category).Error; err != nil {
		return nil, errors.New("failed to update category")
	}

	return &models.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil
}

func (cs *CategoryService) DeleteCategory(categoryID uint) error {
	var category models.Category
	if err := database.DB.First(&category, categoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("category not found")
		}
		return errors.New("database error")
	}

	// Check if category has products
	var count int64
	if err := database.DB.Model(&models.Product{}).Where("category_id = ?", categoryID).Count(&count).Error; err != nil {
		return errors.New("failed to check category products")
	}

	if count > 0 {
		return errors.New("cannot delete category with existing products")
	}

	if err := database.DB.Delete(&category).Error; err != nil {
		return errors.New("failed to delete category")
	}

	return nil
}
