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

// Scheduler CRUD operations
func CreateScheduler(scheduler *models.Scheduler) error {
	return DB.Create(scheduler).Error
}

func GetSchedulersByUserID(userID uint) ([]models.Scheduler, error) {
	var schedulers []models.Scheduler
	err := DB.Where("user_id = ?", userID).Find(&schedulers).Error
	return schedulers, err
}

func GetSchedulerByToken(token string) (models.Scheduler, error) {
	var scheduler models.Scheduler
	err := DB.Where("token = ?", token).Find(&scheduler).Error
	return scheduler, err
}
func GetSchedulersByUserIDWithDevices(userID uint) ([]models.Scheduler, error) {
	var schedulers []models.Scheduler
	err := DB.Where("user_id = ?", userID).Preload("Devices").Find(&schedulers).Error
	return schedulers, err
}

func DeleteScheduler(schedulerID uint) error {
	return DB.Delete(&models.Scheduler{}, schedulerID).Error
}

func UpdateSchedulerTokenLastUsedAt(device *models.Scheduler) error {
	now := time.Now()
	return DB.Model(&models.Scheduler{}).Where("id = ?", device.ID).Update("token_last_used_at", now).Error
}

func IsSchedulerTokenRateLimited(deviceID uint) (bool, error) {
	var device models.Scheduler

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

func GetDeviceDataBySchedulerID(schedulerID uint) ([]models.DeviceData, error) {
	var data []models.DeviceData
	err := DB.Where("scheduler_id = ?", schedulerID).Order("created_at desc").Find(&data).Error
	return data, err
}

// Validation functions
func ValidateScheduler(scheduler *models.Scheduler) error {
	if scheduler.Threshold != 0 && !time.Time(scheduler.OrderDate).IsZero() {
		return errors.New("scheduler cannot have both threshold and date")
	}
	if scheduler.Threshold == 0 && time.Time(scheduler.OrderDate).IsZero() {
		return errors.New("scheduler must have either device or date")
	}

	return nil
}
