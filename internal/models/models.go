package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	TerminalID string      `gorm:"uniqueIndex;not null"`
	Email      string      `gorm:"uniqueIndex;not null"`
	Schedulers []Scheduler `gorm:"foreignKey:UserID"`
}

type DeviceData struct {
	gorm.Model
	SchedulerID uint      `gorm:"not null"`
	Scheduler   Scheduler `gorm:"foreignKey:SchedulerID"`
	Value       float64   `gorm:"not null"`
}

type Scheduler struct {
	gorm.Model
	Name            string         `gorm:"not null"`
	UserID          uint           `gorm:"not null"`
	User            User           `gorm:"foreignKey:UserID"`
	ProductID       string         `gorm:"not null"` // From terminal.shop
	CardID          string         `gorm:"not null"` // From terminal.shop
	AddressID       string         `gorm:"not null"` // From terminal.shop
	Threshold       float64        // For device-based scheduling
	OrderDate       datatypes.Date // For date-based scheduling
	Token           string         `gorm:"not null"`
	TokenLastUsedAt time.Time
}
