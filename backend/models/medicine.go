package models

import (
	"time"

	"gorm.io/gorm"
)

type Medicine struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Name          string         `json:"name" gorm:"not null"`
	Manufacturer  string         `json:"manufacturer"`
	BatchNumber   string         `json:"batch_number"`
	ExpiryDate    Date           `json:"expiry_date" gorm:"type:date"`
	PurchasePrice float64        `json:"purchase_price" gorm:"not null"`
	SellingPrice  float64        `json:"selling_price" gorm:"not null"`
	Quantity      int            `json:"quantity" gorm:"not null;default:0"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}
