package models

import (
	"time"

	"gorm.io/gorm"
)

type Purchase struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	MedicineID   uint           `json:"medicine_id" gorm:"not null"`
	Medicine     Medicine       `json:"medicine" gorm:"foreignKey:MedicineID"`
	Quantity     int            `json:"quantity" gorm:"not null"`
	PurchasePrice float64       `json:"purchase_price" gorm:"not null"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}
