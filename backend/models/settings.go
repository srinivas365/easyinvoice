package models

import (
	"time"

	"gorm.io/gorm"
)

type Settings struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	AccountID    uint           `json:"account_id" gorm:"uniqueIndex;not null"`
	Account      Account        `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CompanyName  string         `json:"company_name" gorm:"not null"`
	Address      string         `json:"address"`
	PhoneNumber  string         `json:"phone_number"`
	LicenseID    string         `json:"license_id" gorm:"not null"`
	GSTINNumber  string         `json:"gstin_number"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}
