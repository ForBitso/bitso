package services

import (
	"errors"

	"go-shop/database"
	"go-shop/models"

	"gorm.io/gorm"
)

type ProductService struct{}

func NewProductService() *ProductService {
	return &ProductService{}
}

func (ps *ProductService) CreateProduct(req *models.ProductCreateRequest) (*models.ProductResponse, error) {
	// Check if category exists
	var category models.Category
	if err := database.DB.First(&category, req.CategoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, errors.New("database error")
	}

	product := models.Product{
		CategoryID:  &req.CategoryID,
		Title:       req.Title,
		Description: req.Description,
		Images:      models.StringArray(req.Images),
		Price:       req.Price,
		Model:       req.Model,
		ExtraInfo:   req.ExtraInfo,
		Stock:       req.Stock,
	}

	if err := database.DB.Create(&product).Error; err != nil {
		return nil, errors.New("failed to create product")
	}

	return &models.ProductResponse{
		ID:          product.ID,
		CategoryID:  product.CategoryID,
		Title:       product.Title,
		Description: product.Description,
		Images:      []string(product.Images),
		Price:       product.Price,
		Model:       product.Model,
		ExtraInfo:   product.ExtraInfo,
		Stock:       product.Stock,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}, nil
}

func (ps *ProductService) GetProducts(categoryID *uint, limit, offset int) ([]models.ProductResponse, error) {
	var products []models.Product
	query := database.DB

	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}

	if err := query.Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, errors.New("failed to get products")
	}

	var productResponses []models.ProductResponse
	for _, product := range products {
		productResponses = append(productResponses, models.ProductResponse{
			ID:          product.ID,
			CategoryID:  product.CategoryID,
			Title:       product.Title,
			Description: product.Description,
			Images:      []string(product.Images),
			Price:       product.Price,
			Model:       product.Model,
			ExtraInfo:   product.ExtraInfo,
			Stock:       product.Stock,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		})
	}

	return productResponses, nil
}

func (ps *ProductService) GetProductByID(productID uint) (*models.ProductResponse, error) {
	var product models.Product
	if err := database.DB.Preload("Category").First(&product, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, errors.New("database error")
	}

	categoryResponse := &models.CategoryResponse{
		ID:          product.Category.ID,
		Name:        product.Category.Name,
		Description: product.Category.Description,
		CreatedAt:   product.Category.CreatedAt,
		UpdatedAt:   product.Category.UpdatedAt,
	}

	return &models.ProductResponse{
		ID:          product.ID,
		CategoryID:  product.CategoryID,
		Title:       product.Title,
		Description: product.Description,
		Images:      []string(product.Images),
		Price:       product.Price,
		Model:       product.Model,
		ExtraInfo:   product.ExtraInfo,
		Stock:       product.Stock,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
		Category:    categoryResponse,
	}, nil
}

func (ps *ProductService) UpdateProduct(productID uint, req *models.ProductUpdateRequest) (*models.ProductResponse, error) {
	var product models.Product
	if err := database.DB.First(&product, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, errors.New("database error")
	}

	// Check if new category exists
	if req.CategoryID != nil {
		var category models.Category
		if err := database.DB.First(&category, *req.CategoryID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("category not found")
			}
			return nil, errors.New("database error")
		}
		product.CategoryID = req.CategoryID
	}

	// Update fields
	if req.Title != "" {
		product.Title = req.Title
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Images != nil {
		product.Images = models.StringArray(req.Images)
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Model != "" {
		product.Model = req.Model
	}
	if req.ExtraInfo != nil {
		product.ExtraInfo = req.ExtraInfo
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}

	if err := database.DB.Save(&product).Error; err != nil {
		return nil, errors.New("failed to update product")
	}

	return &models.ProductResponse{
		ID:          product.ID,
		CategoryID:  product.CategoryID,
		Title:       product.Title,
		Description: product.Description,
		Images:      []string(product.Images),
		Price:       product.Price,
		Model:       product.Model,
		ExtraInfo:   product.ExtraInfo,
		Stock:       product.Stock,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}, nil
}

func (ps *ProductService) DeleteProduct(productID uint) error {
	var product models.Product
	if err := database.DB.First(&product, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return errors.New("database error")
	}

	// Check if product has order items
	var count int64
	if err := database.DB.Model(&models.OrderItem{}).Where("product_id = ?", productID).Count(&count).Error; err != nil {
		return errors.New("failed to check product orders")
	}

	if count > 0 {
		return errors.New("cannot delete product with existing orders")
	}

	if err := database.DB.Delete(&product).Error; err != nil {
		return errors.New("failed to delete product")
	}

	return nil
}

func (ps *ProductService) UpdateStock(productID uint, quantity int) error {
	var product models.Product
	if err := database.DB.First(&product, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return errors.New("database error")
	}

	newStock := product.Stock + quantity
	if newStock < 0 {
		return errors.New("insufficient stock")
	}

	if err := database.DB.Model(&product).Update("stock", newStock).Error; err != nil {
		return errors.New("failed to update stock")
	}

	return nil
}

// SearchProducts searches products with filters and sorting
func (ps *ProductService) SearchProducts(req *models.ProductSearchRequest) ([]models.ProductResponse, int64, error) {
	// Set default values
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	// Build query
	query := database.DB.Model(&models.Product{}).Preload("Category")

	// Apply filters
	if req.Title != "" {
		query = query.Where("title ILIKE ?", "%"+req.Title+"%")
	}

	if req.CategoryID != nil {
		query = query.Where("category_id = ?", *req.CategoryID)
	}

	if req.MinPrice != nil {
		query = query.Where("price >= ?", *req.MinPrice)
	}

	if req.MaxPrice != nil {
		query = query.Where("price <= ?", *req.MaxPrice)
	}

	// Apply sorting
	switch req.SortBy {
	case "price_asc":
		query = query.Order("price ASC")
	case "price_desc":
		query = query.Order("price DESC")
	case "popularity_asc":
		query = query.Order("order_count ASC")
	case "popularity_desc":
		query = query.Order("order_count DESC")
	case "created_at_asc":
		query = query.Order("created_at ASC")
	case "created_at_desc":
		query = query.Order("created_at DESC")
	default:
		// Default sorting by relevance (title match + popularity)
		if req.Title != "" {
			query = query.Order("order_count DESC, title ASC")
		} else {
			query = query.Order("order_count DESC, created_at DESC")
		}
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.New("failed to count products")
	}

	// Get products with pagination
	var products []models.Product
	if err := query.Limit(req.Limit).Offset(req.Offset).Find(&products).Error; err != nil {
		return nil, 0, errors.New("failed to search products")
	}

	// Convert to response format
	var responses []models.ProductResponse
	for _, product := range products {
		response := models.ProductResponse{
			ID:          product.ID,
			CategoryID:  product.CategoryID,
			Title:       product.Title,
			Description: product.Description,
			Images:      []string(product.Images),
			Price:       product.Price,
			Model:       product.Model,
			ExtraInfo:   product.ExtraInfo,
			Stock:       product.Stock,
			OrderCount:  product.OrderCount,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		}

		if product.Category != nil {
			response.Category = &models.CategoryResponse{
				ID:          product.Category.ID,
				Name:        product.Category.Name,
				Description: product.Category.Description,
				CreatedAt:   product.Category.CreatedAt,
				UpdatedAt:   product.Category.UpdatedAt,
			}
		}

		responses = append(responses, response)
	}

	return responses, total, nil
}

// LogSearch logs search queries for analytics
func (ps *ProductService) LogSearch(userID *uint, query string, filters models.JSONB, results int) error {
	searchLog := models.SearchLog{
		UserID:  userID,
		Query:   query,
		Filters: filters,
		Results: results,
	}

	if err := database.DB.Create(&searchLog).Error; err != nil {
		return errors.New("failed to log search")
	}

	return nil
}
