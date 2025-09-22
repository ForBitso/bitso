package routes

import (
	"go-shop/config"
	"go-shop/handlers"
	"go-shop/middleware"
	"go-shop/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(cfg *config.Config) *gin.Engine {
	// Initialize services
	emailService := services.NewEmailService(cfg)
	authService := services.NewAuthService(cfg, emailService)
	userService := services.NewUserService()
	categoryService := services.NewCategoryService()
	productService := services.NewProductService()
	orderService := services.NewOrderService()
	favoriteService := services.NewFavoriteService()
	roleService := services.NewRoleService()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	productHandler := handlers.NewProductHandler(productService)
	orderHandler := handlers.NewOrderHandler(orderService)
	favoriteHandler := handlers.NewFavoriteHandler(favoriteService)
	adminHandler := handlers.NewAdminHandler(categoryService, productService, orderService)
	roleHandler := handlers.NewRoleHandler(roleService)

	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	// Create router
	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Go Shop API is running",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/verify-otp", authHandler.VerifyOTP)
			auth.POST("/login", authHandler.Login)
			auth.POST("/request-password-reset", authHandler.RequestPasswordReset)
			auth.POST("/reset-password", authHandler.ResetPassword)
		}

		// Public routes
		{
			// Category routes (public)
			categories := v1.Group("/categories")
			{
				categories.GET("/", categoryHandler.GetCategories)
				categories.GET("/:id", categoryHandler.GetCategoryByID)
			}

			// Product routes (public)
			products := v1.Group("/products")
			{
				products.GET("/", productHandler.GetProducts)
				products.GET("/search", productHandler.SearchProducts)
				products.GET("/:id", productHandler.GetProductByID)
			}
		}

		// Protected routes (require authentication)
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			// User routes
			user := protected.Group("/user")
			{
				user.GET("/profile", userHandler.GetProfile)
				user.PUT("/profile", userHandler.UpdateProfile)
				user.GET("/:id", userHandler.GetUserByID)
			}

			// Order routes
			orders := protected.Group("/orders")
			{
				orders.POST("/", orderHandler.CreateOrder)
				orders.GET("/", orderHandler.GetUserOrders)
				orders.GET("/:id", orderHandler.GetOrderByID)
				orders.PUT("/:id", orderHandler.UpdateOrderStatus)
				orders.POST("/:id/pay", orderHandler.PayOrder)
				orders.POST("/:id/cancel", orderHandler.CancelOrder)
			}

			// Favorite routes
			favorites := protected.Group("/favorites")
			{
				favorites.POST("/", favoriteHandler.AddToFavorites)
				favorites.GET("/", favoriteHandler.GetUserFavorites)
				favorites.DELETE("/:id", favoriteHandler.RemoveFromFavorites)
				favorites.GET("/check", favoriteHandler.CheckFavorite)
			}
		}

		// Super Admin routes (require super_admin role)
		superAdmin := v1.Group("/super-admin")
		superAdmin.Use(middleware.SuperAdminMiddleware(cfg))
		{
			// Role management
			roles := superAdmin.Group("/roles")
			{
				roles.POST("/assign", roleHandler.AssignRole)
				roles.DELETE("/remove", roleHandler.RemoveRole)
				roles.GET("/", roleHandler.GetAllRoles)
				roles.POST("/", roleHandler.CreateRole)
				roles.GET("/users/:role", roleHandler.GetUsersByRole)
				roles.GET("/user/:id", roleHandler.GetUserRole)
				roles.GET("/all-users", roleHandler.GetAllUsersWithRoles)
			}

			// Full category management (super admin only)
			superAdminCategories := superAdmin.Group("/categories")
			{
				superAdminCategories.POST("/", adminHandler.CreateCategory)
				superAdminCategories.GET("/", adminHandler.GetCategories)
				superAdminCategories.PUT("/:id", adminHandler.UpdateCategory)
				superAdminCategories.DELETE("/:id", adminHandler.DeleteCategory)
			}

			// Full product management (super admin only)
			superAdminProducts := superAdmin.Group("/products")
			{
				superAdminProducts.POST("/", adminHandler.CreateProduct)
				superAdminProducts.GET("/", adminHandler.GetProducts)
				superAdminProducts.PUT("/:id", adminHandler.UpdateProduct)
				superAdminProducts.DELETE("/:id", adminHandler.DeleteProduct)
			}

			// Full order management (super admin only)
			superAdminOrders := superAdmin.Group("/orders")
			{
				superAdminOrders.GET("/", adminHandler.GetAllOrders)
				superAdminOrders.POST("/:id/confirm", adminHandler.ConfirmOrder)
				superAdminOrders.POST("/:id/ship", adminHandler.ShipOrder)
				superAdminOrders.POST("/:id/deliver", adminHandler.DeliverOrder)
				superAdminOrders.POST("/:id/cancel", adminHandler.CancelOrder)
			}
		}

		// Seller routes (require seller or super_admin role)
		seller := v1.Group("/seller")
		seller.Use(middleware.SellerMiddleware(cfg))
		{
			// Product management (sellers can manage products)
			sellerProducts := seller.Group("/products")
			{
				sellerProducts.POST("/", adminHandler.CreateProduct)
				sellerProducts.GET("/", adminHandler.GetProducts)
				sellerProducts.PUT("/:id", adminHandler.UpdateProduct)
				// Sellers cannot delete products
			}

			// Order management (sellers can only ship orders)
			sellerOrders := seller.Group("/orders")
			{
				sellerOrders.GET("/", adminHandler.GetAllOrders)
				sellerOrders.POST("/:id/ship", adminHandler.ShipOrder)
			}
		}
	}
	return router
}
