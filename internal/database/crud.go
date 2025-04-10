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

func GetDeviceByToken(token string) (models.Device, error) {
	var device models.Device
	err := DB.Where("token = ?", token).First(&device).Error
	return device, err
}

func DeleteDevice(deviceID uint) error {
	return DB.Delete(&models.Device{}, deviceID).Error
}

func UpdateDeviceTokenLastUsedAt(device *models.Device) error {
	now := time.Now()
	return DB.Model(&models.User{}).Where("id = ?", device.ID).Update("token_last_used_at", now).Error
}

func IsDeviceTokenRateLimited(deviceID uint) (bool, error) {
	var device models.Device

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
