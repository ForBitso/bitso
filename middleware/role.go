package middleware

import (
	"log"
	"net/http"
	"strings"

	"go-shop/config"
	"go-shop/database"
	"go-shop/models"
	"go-shop/utils"

	"github.com/gin-gonic/gin"
)

// ExtractTokenFromHeader извлекает токен из заголовка Authorization
func ExtractTokenFromHeader(authHeader string) string {
	const bearerPrefix = "Bearer "
	if strings.HasPrefix(authHeader, bearerPrefix) {
		return authHeader[len(bearerPrefix):]
	}
	return ""
}

// RoleMiddleware проверяет, имеет ли пользователь требуемую роль
func RoleMiddleware(requiredRole string, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Authorization header required",
			})
			c.Abort()
			return
		}

		tokenString := ExtractTokenFromHeader(authHeader)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenString, cfg)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Invalid token",
			})
			c.Abort()
			return
		}

		// Получаем пользователя с ролями из базы данных
		var user models.User
		if err := database.DB.Preload("Roles").First(&user, claims.UserID).Error; err != nil {
			log.Printf("RoleMiddleware: Failed to get user %d: %v", claims.UserID, err)
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "User not found",
			})
			c.Abort()
			return
		}

		// Проверяем, имеет ли пользователь требуемую роль
		hasRole := false
		for _, role := range user.Roles {
			if role.Name == requiredRole {
				hasRole = true
				break
			}
		}

		if !hasRole {
			log.Printf("RoleMiddleware: User %d does not have required role %s", claims.UserID, requiredRole)
			c.JSON(http.StatusForbidden, models.ErrorResponse{
				Error: "Access denied: Insufficient permissions",
			})
			c.Abort()
			return
		}

		// Сохраняем информацию о пользователе в контексте
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_roles", user.Roles)
		c.Next()
	}
}

// SuperAdminMiddleware проверяет, является ли пользователь супер-админом
func SuperAdminMiddleware(cfg *config.Config) gin.HandlerFunc {
	return RoleMiddleware(models.ROLE_SUPER_ADMIN, cfg)
}

// SellerMiddleware проверяет, является ли пользователь продавцом или супер-админом
func SellerMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Authorization header required",
			})
			c.Abort()
			return
		}

		tokenString := ExtractTokenFromHeader(authHeader)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenString, cfg)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Invalid token",
			})
			c.Abort()
			return
		}

		// Получаем пользователя с ролями из базы данных
		var user models.User
		if err := database.DB.Preload("Roles").First(&user, claims.UserID).Error; err != nil {
			log.Printf("SellerMiddleware: Failed to get user %d: %v", claims.UserID, err)
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "User not found",
			})
			c.Abort()
			return
		}

		// Проверяем, имеет ли пользователь роль продавца или супер-админа
		hasPermission := false
		for _, role := range user.Roles {
			if role.Name == models.ROLE_SELLER || role.Name == models.ROLE_SUPER_ADMIN {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			log.Printf("SellerMiddleware: User %d does not have seller or super_admin role", claims.UserID)
			c.JSON(http.StatusForbidden, models.ErrorResponse{
				Error: "Access denied: Seller or Super Admin role required",
			})
			c.Abort()
			return
		}

		// Сохраняем информацию о пользователе в контексте
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_roles", user.Roles)
		c.Next()
	}
}

// LogSensitiveOperation логирует чувствительные операции
func LogSensitiveOperation(operation string, userID uint, details string) {
	log.Printf("SENSITIVE_OPERATION: %s | UserID: %d | Details: %s", operation, userID, details)
}
