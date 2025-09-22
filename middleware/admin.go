package middleware

import (
	"log"
	"net/http"

	"go-shop/config"
	"go-shop/database"
	"go-shop/models"
	"go-shop/utils"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware проверяет, является ли пользователь супер-админом или продавцом (для обратной совместимости)
func AdminMiddleware(cfg *config.Config) gin.HandlerFunc {
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
			log.Printf("AdminMiddleware: Failed to get user %d: %v", claims.UserID, err)
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "User not found",
			})
			c.Abort()
			return
		}

		// Проверяем, имеет ли пользователь роль супер-админа или продавца
		hasPermission := false
		for _, role := range user.Roles {
			if role.Name == models.ROLE_SUPER_ADMIN || role.Name == models.ROLE_SELLER {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			log.Printf("AdminMiddleware: User %d does not have admin or seller role", claims.UserID)
			c.JSON(http.StatusForbidden, models.ErrorResponse{
				Error: "Access denied: Admin or Seller role required",
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
