package database

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/softsrv/brewferring/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("brewferring.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	// Auto migrate the schema
	err = DB.AutoMigrate(&models.User{}, &models.Device{}, &models.DeviceToken{}, &models.DeviceData{}, &models.Scheduler{})
	if err != nil {
		return err
	}

	return nil
}

func GenerateDeviceToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func SaveDeviceData(data *models.DeviceData) error {
	return DB.Create(data).Error
}

func GetUserDevices(userID uint) ([]models.Device, error) {
	var devices []models.Device
	err := DB.Where("user_id = ?", userID).Find(&devices).Error
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func GetUserSchedulers(userID uint) ([]models.Scheduler, error) {
	var schedulers []models.Scheduler
	err := DB.Preload("Device").Where("user_id = ?", userID).Find(&schedulers).Error
	if err != nil {
		return nil, err
	}
	return schedulers, nil
}
