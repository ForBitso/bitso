package services

import (
	"errors"
	"log"

	"go-shop/database"
	"go-shop/models"

	"gorm.io/gorm"
)

type RoleService struct{}

func NewRoleService() *RoleService {
	return &RoleService{}
}

// GetUserRole возвращает роль пользователя
func (rs *RoleService) GetUserRole(userID uint) (string, error) {
	var userRole models.UserRole
	err := database.DB.Preload("Role").Where("user_id = ?", userID).First(&userRole).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ROLE_USER, nil // По умолчанию роль user
		}
		return "", err
	}
	return userRole.Role.Name, nil
}

// AssignRole назначает роль пользователю (только для super_admin)
func (rs *RoleService) AssignRole(userID uint, roleName string, assignedBy uint) error {
	// Проверяем, что назначающий имеет роль super_admin
	assignerRole, err := rs.GetUserRole(assignedBy)
	if err != nil {
		return errors.New("failed to verify assigner role")
	}
	if assignerRole != models.ROLE_SUPER_ADMIN {
		return errors.New("only super admin can assign roles")
	}

	// Получаем роль из базы данных
	var role models.Role
	if err := database.DB.Where("name = ?", roleName).First(&role).Error; err != nil {
		return errors.New("role not found")
	}

	// Проверяем, существует ли пользователь
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	// Начинаем транзакцию
	tx := database.DB.Begin()
	if tx.Error != nil {
		return errors.New("failed to start transaction")
	}

	// Удаляем существующую роль пользователя (если есть)
	if err := tx.Where("user_id = ?", userID).Delete(&models.UserRole{}).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to remove existing role")
	}

	// Назначаем новую роль
	userRole := models.UserRole{
		UserID: userID,
		RoleID: role.ID,
	}

	if err := tx.Create(&userRole).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to assign role")
	}

	// Подтверждаем транзакцию
	if err := tx.Commit().Error; err != nil {
		return errors.New("failed to commit role assignment")
	}

	// Логируем операцию
	log.Printf("Role assigned: User %d assigned role %s by user %d", userID, roleName, assignedBy)

	return nil
}

// RemoveRole удаляет роль у пользователя (только для super_admin)
func (rs *RoleService) RemoveRole(userID uint, removedBy uint) error {
	// Проверяем, что удаляющий имеет роль super_admin
	removerRole, err := rs.GetUserRole(removedBy)
	if err != nil {
		return errors.New("failed to verify remover role")
	}
	if removerRole != models.ROLE_SUPER_ADMIN {
		return errors.New("only super admin can remove roles")
	}

	// Проверяем, что пользователь не пытается удалить роль самому себе
	if userID == removedBy {
		return errors.New("cannot remove your own role")
	}

	// Проверяем, существует ли пользователь
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	// Проверяем, есть ли у пользователя роль для удаления
	var userRole models.UserRole
	if err := database.DB.Where("user_id = ?", userID).First(&userRole).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user has no role to remove")
		}
		return errors.New("failed to check user role")
	}

	// Начинаем транзакцию
	tx := database.DB.Begin()
	if tx.Error != nil {
		return errors.New("failed to start transaction")
	}

	// Удаляем роль пользователя
	if err := tx.Where("user_id = ?", userID).Delete(&models.UserRole{}).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to remove role")
	}

	// Подтверждаем транзакцию
	if err := tx.Commit().Error; err != nil {
		return errors.New("failed to commit role removal")
	}

	// Логируем операцию
	log.Printf("Role removed: User %d role removed by user %d", userID, removedBy)

	return nil
}

// GetUsersByRole возвращает всех пользователей с определенной ролью
func (rs *RoleService) GetUsersByRole(roleName string) ([]models.User, error) {
	var users []models.User
	err := database.DB.Preload("Roles").Joins("JOIN user_roles ON users.id = user_roles.user_id").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("roles.name = ?", roleName).Find(&users).Error
	return users, err
}

// GetAllRoles возвращает все роли
func (rs *RoleService) GetAllRoles() ([]models.Role, error) {
	var roles []models.Role
	err := database.DB.Find(&roles).Error
	return roles, err
}

// CreateRole создает новую роль (только для super_admin)
func (rs *RoleService) CreateRole(req *models.RoleCreateRequest, createdBy uint) (*models.Role, error) {
	// Проверяем, что создающий имеет роль super_admin
	creatorRole, err := rs.GetUserRole(createdBy)
	if err != nil {
		return nil, errors.New("failed to verify creator role")
	}
	if creatorRole != models.ROLE_SUPER_ADMIN {
		return nil, errors.New("only super admin can create roles")
	}

	role := models.Role{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := database.DB.Create(&role).Error; err != nil {
		return nil, errors.New("failed to create role")
	}

	// Логируем операцию
	log.Printf("Role created: Role %s created by user %d", req.Name, createdBy)

	return &role, nil
}

// HasRole проверяет, имеет ли пользователь определенную роль
func (rs *RoleService) HasRole(userID uint, roleName string) bool {
	userRole, err := rs.GetUserRole(userID)
	if err != nil {
		return false
	}
	return userRole == roleName
}

// IsSuperAdmin проверяет, является ли пользователь супер-админом
func (rs *RoleService) IsSuperAdmin(userID uint) bool {
	return rs.HasRole(userID, models.ROLE_SUPER_ADMIN)
}

// IsSeller проверяет, является ли пользователь продавцом
func (rs *RoleService) IsSeller(userID uint) bool {
	return rs.HasRole(userID, models.ROLE_SELLER)
}

// GetAllUsersWithRoles возвращает всех пользователей с их ролями
func (rs *RoleService) GetAllUsersWithRoles() ([]models.User, error) {
	var users []models.User
	err := database.DB.Preload("Roles").Find(&users).Error
	return users, err
}

// FixUsersWithoutRoles назначает роль user всем пользователям без ролей
func (rs *RoleService) FixUsersWithoutRoles() error {
	// Получаем роль user
	var userRole models.Role
	if err := database.DB.Where("name = ?", models.ROLE_USER).First(&userRole).Error; err != nil {
		return errors.New("user role not found")
	}

	// Находим всех пользователей без ролей
	var usersWithoutRoles []models.User
	err := database.DB.Raw(`
		SELECT u.* FROM users u 
		LEFT JOIN user_roles ur ON u.id = ur.user_id 
		WHERE ur.user_id IS NULL
	`).Scan(&usersWithoutRoles).Error

	if err != nil {
		return err
	}

	// Назначаем роль user каждому пользователю без роли
	for _, user := range usersWithoutRoles {
		userRoleAssignment := models.UserRole{
			UserID: user.ID,
			RoleID: userRole.ID,
		}
		if err := database.DB.Create(&userRoleAssignment).Error; err != nil {
			log.Printf("Failed to assign role to user %d: %v", user.ID, err)
		} else {
			log.Printf("Assigned default role to user: %s (ID: %d)", user.Email, user.ID)
		}
	}

	log.Printf("Fixed roles for %d users", len(usersWithoutRoles))
	return nil
}
