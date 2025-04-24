package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	TerminalID string      `gorm:"uniqueIndex;not null"`
	Email      string      `gorm:"uniqueIndex;not null"`
	Devices    []Device    `gorm:"foreignKey:UserID"`
	Schedulers []Scheduler `gorm:"foreignKey:UserID"`
}

type Device struct {
	gorm.Model
	Name            string `gorm:"not null"`
	UserID          uint   `gorm:"not null"`
	User            User   `gorm:"foreignKey:UserID"`
	Token           string `gorm:"not null"`
	TokenLastUsedAt time.Time
	Data            []DeviceData `gorm:"foreignKey:DeviceID"`
}

type DeviceData struct {
	gorm.Model
	DeviceID uint    `gorm:"not null"`
	Device   Device  `gorm:"foreignKey:DeviceID"`
	Value    float64 `gorm:"not null"`
}

type Scheduler struct {
	gorm.Model
	Name      string    `gorm:"not null"`
	UserID    uint      `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID"`
	DeviceID  uint      `gorm:"not null"`
	Device    Device    `gorm:"foreignKey:DeviceID"`
	ProductID string    `gorm:"not null"` // From terminal.shop
	Threshold float64   // For device-based scheduling
	Date      time.Time // For date-based scheduling

}
