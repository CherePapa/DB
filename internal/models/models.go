package models

import (
	"time"

	"gorm.io/gorm"
)

type Medicine struct {
	gorm.Model
	Name         string    `gorm:"size:255;not null"`
	Manufacturer string    `gorm:"size:255;not null"`
	Price        float64   `gorm:"type:decimal(10,2);not null"`
	Quantity     int       `gorm:"not null"`
	BatchNumber  string    `gorm:"size:255;not null"`
	ExpiryDate   time.Time `gorm:"not null"`
}
