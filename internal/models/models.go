package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	TerminalID string      `gorm:"uniqueIndex;not null"`
	Email      string      `gorm:"uniqueIndex;not null"`
	Name       string      `gorm:"not null"`
	CreatedAt  time.Time   `gorm:"not null"`
	Devices    []Device    `gorm:"foreignKey:UserID"`
	Schedulers []Scheduler `gorm:"foreignKey:UserID"`
}

type Device struct {
	gorm.Model
	Name       string        `gorm:"not null"`
	UserID     uint          `gorm:"not null"`
	User       User          `gorm:"foreignKey:UserID"`
	Tokens     []DeviceToken `gorm:"foreignKey:DeviceID"`
	Data       []DeviceData  `gorm:"foreignKey:DeviceID"`
	Schedulers []Scheduler   `gorm:"foreignKey:DeviceID"`
}

type DeviceToken struct {
	gorm.Model
	Token      string `gorm:"uniqueIndex;not null"`
	DeviceID   uint   `gorm:"not null"`
	Device     Device `gorm:"foreignKey:DeviceID"`
	LastUsedAt time.Time
}

type DeviceData struct {
	gorm.Model
	DeviceID      uint        `gorm:"not null"`
	Device        Device      `gorm:"foreignKey:DeviceID"`
	DeviceTokenID uint        `gorm:"not null"`
	DeviceToken   DeviceToken `gorm:"foreignKey:DeviceTokenID"`
	Value         float64     `gorm:"not null"`
	CreatedAt     time.Time   `gorm:"not null"`
}

type Scheduler struct {
	gorm.Model
	Name      string    `gorm:"not null"`
	UserID    uint      `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID"`
	DeviceID  uint      `gorm:"not null"`
	Device    Device    `gorm:"foreignKey:DeviceID"`
	Type      string    `gorm:"not null"` // "device" or "date"
	Threshold float64   // For device-based scheduling
	Date      time.Time // For date-based scheduling
}
