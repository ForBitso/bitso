package handlers

import (
	"net/http"
	"strconv"

	"go-shop/models"
	"go-shop/services"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	roleService *services.RoleService
}

func NewRoleHandler(roleService *services.RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

// AssignRole назначает роль пользователю (только для super_admin)
func (rh *RoleHandler) AssignRole(c *gin.Context) {
	var req models.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request data",
			Message: err.Error(),
		})
		return
	}

	// Получаем ID текущего пользователя из контекста
	currentUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	err := rh.roleService.AssignRole(req.UserID, req.Role, currentUserID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to assign role",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Role assigned successfully",
	})
}

// RemoveRole удаляет роль у пользователя (только для super_admin)
func (rh *RoleHandler) RemoveRole(c *gin.Context) {
	var req models.RemoveRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request data",
			Message: err.Error(),
		})
		return
	}

	// Получаем ID текущего пользователя из контекста
	currentUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	err := rh.roleService.RemoveRole(req.UserID, currentUserID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to remove role",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Role removed successfully",
	})
}

// GetUsersByRole возвращает всех пользователей с определенной ролью
func (rh *RoleHandler) GetUsersByRole(c *gin.Context) {
	roleName := c.Param("role")
	if roleName == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Role name is required",
		})
		return
	}

	users, err := rh.roleService.GetUsersByRole(roleName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get users",
			Message: err.Error(),
		})
		return
	}

	// Конвертируем в response format
	var userResponses []models.UserResponse
	for _, user := range users {
		var roleResponses []models.RoleResponse
		for _, role := range user.Roles {
			roleResponses = append(roleResponses, models.RoleResponse{
				ID:          role.ID,
				Name:        role.Name,
				Description: role.Description,
				CreatedAt:   role.CreatedAt,
				UpdatedAt:   role.UpdatedAt,
			})
		}

		userResponses = append(userResponses, models.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Roles:     roleResponses,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Users retrieved successfully",
		Data:    userResponses,
	})
}

// GetAllRoles возвращает все роли
func (rh *RoleHandler) GetAllRoles(c *gin.Context) {
	roles, err := rh.roleService.GetAllRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get roles",
			Message: err.Error(),
		})
		return
	}

	// Конвертируем в response format
	var roleResponses []models.RoleResponse
	for _, role := range roles {
		roleResponses = append(roleResponses, models.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Roles retrieved successfully",
		Data:    roleResponses,
	})
}

// CreateRole создает новую роль (только для super_admin)
func (rh *RoleHandler) CreateRole(c *gin.Context) {
	var req models.RoleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request data",
			Message: err.Error(),
		})
		return
	}

	// Получаем ID текущего пользователя из контекста
	currentUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	role, err := rh.roleService.CreateRole(&req, currentUserID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to create role",
			Message: err.Error(),
		})
		return
	}

	roleResponse := models.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Message: "Role created successfully",
		Data:    roleResponse,
	})
}

// GetUserRole возвращает роль конкретного пользователя
func (rh *RoleHandler) GetUserRole(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	role, err := rh.roleService.GetUserRole(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get user role",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "User role retrieved successfully",
		Data: map[string]string{
			"role": role,
		},
	})
}

// GetAllUsersWithRoles возвращает всех пользователей с их ролями
func (rh *RoleHandler) GetAllUsersWithRoles(c *gin.Context) {
	users, err := rh.roleService.GetAllUsersWithRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get users",
			Message: err.Error(),
		})
		return
	}

	// Конвертируем в response format
	var userResponses []models.UserResponse
	for _, user := range users {
		var roleResponses []models.RoleResponse
		for _, role := range user.Roles {
			roleResponses = append(roleResponses, models.RoleResponse{
				ID:          role.ID,
				Name:        role.Name,
				Description: role.Description,
				CreatedAt:   role.CreatedAt,
				UpdatedAt:   role.UpdatedAt,
			})
		}

		userResponses = append(userResponses, models.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Roles:     roleResponses,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Users retrieved successfully",
		Data:    userResponses,
	})
}
