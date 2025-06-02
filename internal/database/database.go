package database

import (
	"crypto/rand"
	"encoding/hex"

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
	err = DB.AutoMigrate(&models.User{}, &models.DeviceData{}, &models.Buffer{})
	if err != nil {
		return err
	}

	return nil
}

func GenerateDeviceToken() (string, error) {
	b := make([]byte, 20)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return "dt_" + hex.EncodeToString(b), nil
}

func SaveDeviceData(data *models.DeviceData) error {
	return DB.Create(data).Error
}
