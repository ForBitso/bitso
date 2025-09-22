package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"go-shop/config"
	"go-shop/database"
	"go-shop/models"
	"go-shop/utils"

	"gorm.io/gorm"
)

type AuthService struct {
	config       *config.Config
	emailService *EmailService
}

func NewAuthService(cfg *config.Config, emailService *EmailService) *AuthService {
	return &AuthService{
		config:       cfg,
		emailService: emailService,
	}
}

func (as *AuthService) Register(req *models.UserCreateRequest) (*models.UserResponse, error) {
	// Check if user already exists
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("user with this email already exists")
	}

	// Check if user is pending verification
	ctx := context.Background()
	if database.CheckPendingUserExists(ctx, req.Email) {
		return nil, errors.New("user registration is pending verification")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Generate OTP
	otp, err := utils.GenerateOTP(as.config.OTP.Length)
	if err != nil {
		return nil, errors.New("failed to generate OTP")
	}

	// Store pending user data in Redis
	pendingUser := models.UserCreateRequest{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	expiration := time.Duration(as.config.OTP.ExpireMinutes) * time.Minute
	if err := database.SetPendingUser(ctx, req.Email, pendingUser, expiration); err != nil {
		return nil, errors.New("failed to store pending user data")
	}

	// Store OTP in Redis
	if err := database.SetOTP(ctx, req.Email, otp, expiration); err != nil {
		return nil, errors.New("failed to store OTP")
	}

	// Send OTP email
	if err := as.emailService.SendOTPEmail(req.Email, otp); err != nil {
		log.Printf("Failed to send OTP email: %v", err)
		return nil, fmt.Errorf("failed to send OTP email: %v", err)
	}

	return &models.UserResponse{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  false,
	}, nil
}

func (as *AuthService) VerifyOTP(req *models.OTPVerifyRequest) (*models.UserResponse, error) {
	ctx := context.Background()

	// Verify OTP
	storedOTP, err := database.GetOTP(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid or expired OTP")
	}

	if storedOTP != req.OTP {
		return nil, errors.New("invalid OTP")
	}

	// Get pending user data
	pendingUserJSON, err := database.GetPendingUser(ctx, req.Email)
	if err != nil {
		return nil, errors.New("pending user data not found")
	}

	var pendingUser models.UserCreateRequest
	if err := json.Unmarshal([]byte(pendingUserJSON), &pendingUser); err != nil {
		return nil, errors.New("invalid pending user data")
	}

	// Create user in database
	user := models.User{
		Email:     pendingUser.Email,
		Password:  pendingUser.Password,
		FirstName: pendingUser.FirstName,
		LastName:  pendingUser.LastName,
		IsActive:  true,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return nil, errors.New("failed to create user")
	}

	// Assign default user role
	var userRole models.Role
	if err := database.DB.Where("name = ?", models.ROLE_USER).First(&userRole).Error; err != nil {
		log.Printf("Warning: Failed to find user role: %v", err)
	} else {
		// Create user role assignment
		userRoleAssignment := models.UserRole{
			UserID: user.ID,
			RoleID: userRole.ID,
		}
		if err := database.DB.Create(&userRoleAssignment).Error; err != nil {
			log.Printf("Warning: Failed to assign default role to user: %v", err)
		} else {
			log.Printf("Assigned default role 'user' to new user: %s", user.Email)
		}
	}

	// Clean up Redis data
	database.DeleteOTP(ctx, req.Email)
	database.DeletePendingUser(ctx, req.Email)

	// Send welcome email
	if err := as.emailService.SendWelcomeEmail(user.Email, user.FirstName); err != nil {
		log.Printf("Failed to send welcome email: %v", err)
	}

	// Get user with roles for response
	var userWithRoles models.User
	if err := database.DB.Preload("Roles").First(&userWithRoles, user.ID).Error; err != nil {
		log.Printf("Warning: Failed to load user roles: %v", err)
		userWithRoles = user // Use user without roles if loading fails
	}

	// Convert roles to response format
	var roleResponses []models.RoleResponse
	for _, role := range userWithRoles.Roles {
		roleResponses = append(roleResponses, models.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
		})
	}

	return &models.UserResponse{
		ID:        userWithRoles.ID,
		Email:     userWithRoles.Email,
		FirstName: userWithRoles.FirstName,
		LastName:  userWithRoles.LastName,
		Roles:     roleResponses,
		IsActive:  userWithRoles.IsActive,
		CreatedAt: userWithRoles.CreatedAt,
	}, nil
}

func (as *AuthService) Login(req *models.UserLoginRequest) (*models.LoginResponse, error) {
	// Find user with roles
	var user models.User
	if err := database.DB.Preload("Roles").Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, errors.New("database error")
	}

	if !user.IsActive {
		return nil, errors.New("account is not activated")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	// Get user's primary role (first role or default to user)
	userRole := models.ROLE_USER
	if len(user.Roles) > 0 {
		userRole = user.Roles[0].Name
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, userRole, as.config)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// Convert roles to response format
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

	response := &models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Roles:     roleResponses,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}

	// Return response with token
	return &models.LoginResponse{
		User:  *response,
		Token: token,
	}, nil
}

func (as *AuthService) RequestPasswordReset(req *models.PasswordResetRequest) error {
	// Check if user exists
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Don't reveal if user exists or not
			return nil
		}
		return errors.New("database error")
	}

	// Generate OTP
	otp, err := utils.GenerateOTP(as.config.OTP.Length)
	if err != nil {
		return errors.New("failed to generate OTP")
	}

	// Store OTP in Redis
	ctx := context.Background()
	expiration := time.Duration(as.config.OTP.ExpireMinutes) * time.Minute
	if err := database.SetPasswordResetToken(ctx, req.Email, otp, expiration); err != nil {
		return errors.New("failed to store reset token")
	}

	// Send password reset email
	if err := as.emailService.SendPasswordResetEmail(req.Email, otp); err != nil {
		log.Printf("Failed to send password reset email: %v", err)
		return fmt.Errorf("failed to send password reset email: %v", err)
	}

	return nil
}

func (as *AuthService) ResetPassword(req *models.PasswordResetConfirmRequest) error {
	ctx := context.Background()

	// Verify OTP
	storedOTP, err := database.GetPasswordResetToken(ctx, req.Email)
	if err != nil {
		return errors.New("invalid or expired reset token")
	}

	if storedOTP != req.OTP {
		return errors.New("invalid reset token")
	}

	// Find user
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return errors.New("user not found")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// Update password
	if err := database.DB.Model(&user).Update("password", hashedPassword).Error; err != nil {
		return errors.New("failed to update password")
	}

	// Clean up Redis data
	database.DeletePasswordResetToken(ctx, req.Email)

	return nil
}
