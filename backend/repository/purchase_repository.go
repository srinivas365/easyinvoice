package repository

import (
	"pharmacy-erp/config"
	"pharmacy-erp/models"
)

type PurchaseRepository struct{}

func NewPurchaseRepository() *PurchaseRepository {
	return &PurchaseRepository{}
}

func (r *PurchaseRepository) Create(purchase *models.Purchase) error {
	return config.DB.Create(purchase).Error
}

func (r *PurchaseRepository) FindAll() ([]models.Purchase, error) {
	var purchases []models.Purchase
	err := config.DB.Preload("Medicine").Order("created_at DESC").Find(&purchases).Error
	return purchases, err
}

func (r *PurchaseRepository) FindByID(id uint) (*models.Purchase, error) {
	var purchase models.Purchase
	err := config.DB.Preload("Medicine").First(&purchase, id).Error
	if err != nil {
		return nil, err
	}
	return &purchase, nil
}
