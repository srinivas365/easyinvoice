package models

import (
	"time"

	"gorm.io/gorm"
)

type PurchaseInvoice struct {
	ID                uint                    `json:"id" gorm:"primaryKey"`
	DistributorID     uint                    `json:"distributor_id" gorm:"not null"`
	Distributor       Distributor             `json:"distributor" gorm:"foreignKey:DistributorID"`
	PurchaseDate      time.Time               `json:"purchase_date" gorm:"not null"`
	AdditionalDiscount float64                `json:"additional_discount" gorm:"default:0"`
	TotalAmount       float64                 `json:"total_amount" gorm:"not null"`
	Items             []PurchaseInvoiceItem   `json:"items" gorm:"foreignKey:PurchaseInvoiceID"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
	DeletedAt         gorm.DeletedAt         `json:"-" gorm:"index"`
}

type PurchaseInvoiceItem struct {
	ID                 uint           `json:"id" gorm:"primaryKey"`
	PurchaseInvoiceID  uint           `json:"purchase_invoice_id" gorm:"not null"`
	ProductID          uint           `json:"product_id" gorm:"not null;index"`
	Product            *Product       `json:"product" gorm:"foreignKey:ProductID"`
	BatchNo            string         `json:"batch_no"`
	Expiry             time.Time      `json:"expiry"`
	Quantity           float64        `json:"quantity" gorm:"not null"`
	Free               float64        `json:"free" gorm:"default:0"`
	Mrp                float64        `json:"mrp" gorm:"default:0"`
	Rate               float64        `json:"rate" gorm:"not null"`
	Amount             float64        `json:"amount" gorm:"not null"`
	Gst                float64        `json:"gst" gorm:"default:0"`
	DiscountPercent    float64        `json:"discount_percent" gorm:"default:0"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt  `json:"-" gorm:"index"`
}
