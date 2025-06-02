package database

import (
	"errors"
	"time"

	"github.com/softsrv/brewferring/internal/models"
)

// User CRUD operations
func CreateUser(user *models.User) error {
	return DB.Create(user).Error
}

func GetUserByTerminalID(terminalID string) (*models.User, error) {
	var user models.User
	err := DB.Where("terminal_id = ?", terminalID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Buffer CRUD operations
func CreateBuffer(buffer *models.Buffer) error {
	return DB.Create(buffer).Error
}

func GetBuffersByUserID(userID uint) ([]models.Buffer, error) {
	var buffers []models.Buffer
	err := DB.Where("user_id = ?", userID).Find(&buffers).Error
	return buffers, err
}

func GetBufferByToken(token string) (models.Buffer, error) {
	var buffer models.Buffer
	err := DB.Where("token = ?", token).Find(&buffer).Error
	return buffer, err
}
func GetBuffersByUserIDWithDevices(userID uint) ([]models.Buffer, error) {
	var buffers []models.Buffer
	err := DB.Where("user_id = ?", userID).Preload("Devices").Find(&buffers).Error
	return buffers, err
}

func DeleteBuffer(bufferID uint) error {
	return DB.Delete(&models.Buffer{}, bufferID).Error
}

func UpdateBufferTokenLastUsedAt(device *models.Buffer) error {
	now := time.Now()
	return DB.Model(&models.Buffer{}).Where("id = ?", device.ID).Update("token_last_used_at", now).Error
}

func IsBufferTokenRateLimited(deviceID uint) (bool, error) {
	var device models.Buffer

	err := DB.Where("id = ?", deviceID).First(&device).Error
	if err != nil {
		return true, err
	}

	// Check if more than 1 hour has passed since last use
	return CheckIsTokenLimited(device.TokenLastUsedAt)
}

func CheckIsTokenLimited(lastUsed time.Time) (bool, error) {
	if lastUsed.IsZero() {
		return false, nil
	}

	// Check if more than 1 hour has passed since last use
	return time.Since(lastUsed) < time.Hour, nil
}

// DeviceData CRUD operations
func CreateDeviceData(data *models.DeviceData) error {
	return DB.Create(data).Error
}

func GetDeviceDataByBufferID(bufferID uint) ([]models.DeviceData, error) {
	var data []models.DeviceData
	err := DB.Where("buffer_id = ?", bufferID).Order("created_at desc").Find(&data).Error
	return data, err
}

// Validation functions
func ValidateBuffer(buffer *models.Buffer) error {
	if buffer.Threshold != 0 && !time.Time(buffer.OrderDate).IsZero() {
		return errors.New("buffer cannot have both threshold and date")
	}
	if buffer.Threshold == 0 && time.Time(buffer.OrderDate).IsZero() {
		return errors.New("buffer must have either device or date")
	}

	return nil
}
