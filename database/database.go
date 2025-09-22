package database

import (
	"fmt"
	"log"

	"go-shop/config"
	"go-shop/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB(cfg *config.Config) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")
}

func Migrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.UserRole{},
		&models.Category{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
		&models.Favorite{},
		&models.SearchLog{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Создаем базовые роли, если их нет
	createDefaultRoles()

	log.Println("Database migration completed")
}

// createDefaultRoles создает базовые роли в системе
func createDefaultRoles() {
	roles := []models.Role{
		{Name: models.ROLE_SUPER_ADMIN, Description: "Super Administrator with full access"},
		{Name: models.ROLE_SELLER, Description: "Seller with limited admin access"},
		{Name: models.ROLE_USER, Description: "Regular user"},
	}

	for _, role := range roles {
		var existingRole models.Role
		if err := DB.Where("name = ?", role.Name).First(&existingRole).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				DB.Create(&role)
				log.Printf("Created default role: %s", role.Name)
			}
		}
	}
}

func GetDB() *gorm.DB {
	return DB
}
