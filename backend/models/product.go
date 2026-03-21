package models

import (
	"time"

	"gorm.io/gorm"
)

// Product is the catalog of medicines/products. Purchase invoice items reference this.
// The Medicines section lists Products with aggregated quantity from purchase invoice items.
type Product struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	Mfg       string         `json:"mfg"`
	HsnCode   string         `json:"hsn_code"`
	Pack      string         `json:"pack"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
