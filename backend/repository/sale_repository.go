package repository

import (
	"pharmacy-erp/config"
	"pharmacy-erp/models"
	"time"
)

type SaleRepository struct{}

func NewSaleRepository() *SaleRepository {
	return &SaleRepository{}
}

func (r *SaleRepository) Create(sale *models.Sale) error {
	return config.DB.Create(sale).Error
}

func (r *SaleRepository) FindAll() ([]models.Sale, error) {
	var sales []models.Sale
	err := config.DB.Preload("Items.Medicine").Where("deleted_at IS NULL").Order("created_at DESC").Find(&sales).Error
	return sales, err
}

func (r *SaleRepository) FindAllArchived() ([]models.Sale, error) {
	var sales []models.Sale
	err := config.DB.Unscoped().Preload("Items.Medicine").Where("deleted_at IS NOT NULL").Order("deleted_at DESC").Find(&sales).Error
	return sales, err
}

func (r *SaleRepository) FindByID(id uint) (*models.Sale, error) {
	var sale models.Sale
	err := config.DB.Preload("Items.Medicine").Where("deleted_at IS NULL").First(&sale, id).Error
	if err != nil {
		return nil, err
	}
	return &sale, nil
}

func (r *SaleRepository) FindByIDUnscoped(id uint) (*models.Sale, error) {
	var sale models.Sale
	err := config.DB.Unscoped().Preload("Items.Medicine").First(&sale, id).Error
	if err != nil {
		return nil, err
	}
	return &sale, nil
}

func (r *SaleRepository) GetTodaySales() (float64, error) {
	var total float64
	today := time.Now().Format("2006-01-02")
	err := config.DB.Model(&models.Sale{}).
		Where("DATE(created_at) = ?", today).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&total).Error
	return total, err
}

func (r *SaleRepository) GetTotalRevenue() (float64, error) {
	var total float64
	err := config.DB.Model(&models.Sale{}).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&total).Error
	return total, err
}

func (r *SaleRepository) Update(sale *models.Sale) error {
	return config.DB.Save(sale).Error
}

func (r *SaleRepository) Delete(id uint) error {
	// Soft delete sale items first
	err := config.DB.Where("sale_id = ?", id).Delete(&models.SaleItem{}).Error
	if err != nil {
		return err
	}
	// Soft delete sale
	return config.DB.Delete(&models.Sale{}, id).Error
}

func (r *SaleRepository) Restore(id uint) error {
	// Restore sale
	err := config.DB.Unscoped().Model(&models.Sale{}).Where("id = ?", id).Update("deleted_at", nil).Error
	if err != nil {
		return err
	}
	// Restore sale items
	return config.DB.Unscoped().Model(&models.SaleItem{}).Where("sale_id = ?", id).Update("deleted_at", nil).Error
}
