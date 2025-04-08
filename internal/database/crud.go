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

// Device CRUD operations
func CreateDevice(device *models.Device) error {
	return DB.Create(device).Error
}

func GetDevicesByUserID(userID uint) ([]models.Device, error) {
	var devices []models.Device
	err := DB.Where("user_id = ?", userID).Find(&devices).Error
	return devices, err
}

func DeleteDevice(deviceID uint) error {
	return DB.Delete(&models.Device{}, deviceID).Error
}

func GetDeviceToken(token string) (*models.DeviceToken, error) {
	var deviceToken models.DeviceToken
	err := DB.Where("token = ?", token).First(&deviceToken).Error
	if err != nil {
		return nil, err
	}
	return &deviceToken, nil
}

func UpdateDeviceTokenLastUsedAt(tokenID uint) error {
	now := time.Now()
	return DB.Model(&models.DeviceToken{}).Where("id = ?", tokenID).Update("last_used_at", now).Error
}

func IsDeviceTokenRateLimited(tokenID uint) (bool, error) {
	var token models.DeviceToken
	err := DB.First(&token, tokenID).Error
	if err != nil {
		return true, err
	}

	if token.LastUsedAt.IsZero() {
		return false, nil
	}

	// Check if more than 1 hour has passed since last use
	return time.Since(token.LastUsedAt) < time.Hour, nil
}

// Scheduler CRUD operations
func CreateScheduler(scheduler *models.Scheduler) error {
	return DB.Create(scheduler).Error
}

func GetSchedulersByUserID(userID uint) ([]models.Scheduler, error) {
	var schedulers []models.Scheduler
	err := DB.Where("user_id = ?", userID).Find(&schedulers).Error
	return schedulers, err
}

func DeleteScheduler(schedulerID uint) error {
	return DB.Delete(&models.Scheduler{}, schedulerID).Error
}

// DeviceData CRUD operations
func CreateDeviceData(data *models.DeviceData) error {
	return DB.Create(data).Error
}

func GetDeviceDataByDeviceID(deviceID uint) ([]models.DeviceData, error) {
	var data []models.DeviceData
	err := DB.Where("device_id = ?", deviceID).Order("created_at desc").Find(&data).Error
	return data, err
}

// Validation functions
func ValidateScheduler(scheduler *models.Scheduler) error {
	if scheduler.DeviceID != 0 && !scheduler.Date.IsZero() {
		return errors.New("scheduler cannot have both device and date")
	}
	if scheduler.DeviceID == 0 && scheduler.Date.IsZero() {
		return errors.New("scheduler must have either device or date")
	}
	if scheduler.DeviceID != 0 && scheduler.Threshold == 0 {
		return errors.New("scheduler with device must have threshold")
	}
	return nil
}
