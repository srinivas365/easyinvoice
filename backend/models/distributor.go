package models

import (
	"time"

	"gorm.io/gorm"
)

type Distributor struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Name         string         `json:"name" gorm:"not null"`
	ContactPerson string        `json:"contact_person"`
	Phone        string         `json:"phone"`
	Email        string         `json:"email"`
	Address      string         `json:"address"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}
