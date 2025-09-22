package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go-shop/config"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func ConnectRedis(cfg *config.Config) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	log.Println("Redis connected successfully")
}

func GetRedis() *redis.Client {
	return RedisClient
}

// SetOTP stores OTP in Redis with expiration
func SetOTP(ctx context.Context, email, otp string, expiration time.Duration) error {
	key := fmt.Sprintf("otp:%s", email)
	return RedisClient.Set(ctx, key, otp, expiration).Err()
}

// GetOTP retrieves OTP from Redis
func GetOTP(ctx context.Context, email string) (string, error) {
	key := fmt.Sprintf("otp:%s", email)
	return RedisClient.Get(ctx, key).Result()
}

// DeleteOTP removes OTP from Redis
func DeleteOTP(ctx context.Context, email string) error {
	key := fmt.Sprintf("otp:%s", email)
	return RedisClient.Del(ctx, key).Err()
}

// SetPasswordResetToken stores password reset token in Redis
func SetPasswordResetToken(ctx context.Context, email, token string, expiration time.Duration) error {
	key := fmt.Sprintf("password_reset:%s", email)
	return RedisClient.Set(ctx, key, token, expiration).Err()
}

// GetPasswordResetToken retrieves password reset token from Redis
func GetPasswordResetToken(ctx context.Context, email string) (string, error) {
	key := fmt.Sprintf("password_reset:%s", email)
	return RedisClient.Get(ctx, key).Result()
}

// DeletePasswordResetToken removes password reset token from Redis
func DeletePasswordResetToken(ctx context.Context, email string) error {
	key := fmt.Sprintf("password_reset:%s", email)
	return RedisClient.Del(ctx, key).Err()
}

// SetPendingUser stores pending user data in Redis
func SetPendingUser(ctx context.Context, email string, userData interface{}, expiration time.Duration) error {
	key := fmt.Sprintf("pending_user:%s", email)
	userJSON, err := json.Marshal(userData)
	if err != nil {
		return err
	}
	return RedisClient.Set(ctx, key, userJSON, expiration).Err()
}

// GetPendingUser retrieves pending user data from Redis
func GetPendingUser(ctx context.Context, email string) (string, error) {
	key := fmt.Sprintf("pending_user:%s", email)
	return RedisClient.Get(ctx, key).Result()
}

// DeletePendingUser removes pending user data from Redis
func DeletePendingUser(ctx context.Context, email string) error {
	key := fmt.Sprintf("pending_user:%s", email)
	return RedisClient.Del(ctx, key).Err()
}

// CheckPendingUserExists checks if pending user exists in Redis
func CheckPendingUserExists(ctx context.Context, email string) bool {
	key := fmt.Sprintf("pending_user:%s", email)
	_, err := RedisClient.Get(ctx, key).Result()
	return err == nil
}
