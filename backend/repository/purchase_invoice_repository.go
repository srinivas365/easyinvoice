package repository

import (
	"pharmacy-erp/config"
	"pharmacy-erp/models"
)

type PurchaseInvoiceRepository struct{}

func NewPurchaseInvoiceRepository() *PurchaseInvoiceRepository {
	return &PurchaseInvoiceRepository{}
}

func (r *PurchaseInvoiceRepository) Create(inv *models.PurchaseInvoice) error {
	return config.DB.Create(inv).Error
}

func (r *PurchaseInvoiceRepository) FindAll() ([]models.PurchaseInvoice, error) {
	var list []models.PurchaseInvoice
	err := config.DB.Preload("Items.Product").Preload("Distributor").Order("purchase_date DESC, created_at DESC").Find(&list).Error
	return list, err
}

func (r *PurchaseInvoiceRepository) FindByID(id uint) (*models.PurchaseInvoice, error) {
	var inv models.PurchaseInvoice
	err := config.DB.Preload("Items.Product").Preload("Distributor").First(&inv, id).Error
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *PurchaseInvoiceRepository) Update(inv *models.PurchaseInvoice) error {
	return config.DB.Save(inv).Error
}

func (r *PurchaseInvoiceRepository) Delete(id uint) error {
	return config.DB.Delete(&models.PurchaseInvoice{}, id).Error
}

func (r *PurchaseInvoiceRepository) DeleteItemsByInvoiceID(invoiceID uint) error {
	return config.DB.Unscoped().Where("purchase_invoice_id = ?", invoiceID).Delete(&models.PurchaseInvoiceItem{}).Error
}
