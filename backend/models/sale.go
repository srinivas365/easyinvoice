package models

import (
	"time"

	"gorm.io/gorm"
)

type Sale struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	CustomerName  string         `json:"customer_name"`
	CustomerNumber string        `json:"customer_number"`
	PrescribedBy  string         `json:"prescribed_by"`
	TotalAmount   float64       `json:"total_amount" gorm:"not null"`
	Items         []SaleItem     `json:"items" gorm:"foreignKey:SaleID"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

type SaleItem struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	SaleID       uint           `json:"sale_id" gorm:"not null"`
	MedicineID   *uint          `json:"medicine_id" gorm:"index;default:null"` // Nullable for manual entries
	Medicine     *Medicine      `json:"medicine" gorm:"foreignKey:MedicineID"`
	MedicineName string         `json:"medicine_name" gorm:"not null"` // Particulars for manual entries
	SNo          int            `json:"s_no" gorm:"default:0"`
	Pack         string         `json:"pack"`
	Quantity     int            `json:"quantity" gorm:"not null"`
	Batch        string         `json:"batch"`
	Expiry       *time.Time     `json:"expiry"`
	MRP          float64        `json:"mrp" gorm:"default:0"`
	Price        float64        `json:"price" gorm:"not null"`
	Discount     float64        `json:"discount" gorm:"default:0"`
	DiscountType string         `json:"discount_type" gorm:"default:'amount'"` // 'amount' or 'percentage'
	Total        float64        `json:"total" gorm:"not null"` // Amount
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}
