package repository

import (
	"pharmacy-erp/config"
	"pharmacy-erp/models"
)

type MedicineRepository struct{}

func NewMedicineRepository() *MedicineRepository {
	return &MedicineRepository{}
}

func (r *MedicineRepository) Create(medicine *models.Medicine) error {
	return config.DB.Create(medicine).Error
}

func (r *MedicineRepository) FindAll(search string) ([]models.Medicine, error) {
	var medicines []models.Medicine
	query := config.DB.Model(&models.Medicine{})
	
	if search != "" {
		query = query.Where("name LIKE ? OR manufacturer LIKE ? OR batch_number LIKE ?", 
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	
	err := query.Find(&medicines).Error
	return medicines, err
}

func (r *MedicineRepository) FindByID(id uint) (*models.Medicine, error) {
	var medicine models.Medicine
	err := config.DB.First(&medicine, id).Error
	if err != nil {
		return nil, err
	}
	return &medicine, nil
}

func (r *MedicineRepository) Update(medicine *models.Medicine) error {
	return config.DB.Save(medicine).Error
}

func (r *MedicineRepository) Delete(id uint) error {
	return config.DB.Delete(&models.Medicine{}, id).Error
}

func (r *MedicineRepository) UpdateStock(id uint, quantity int) error {
	var medicine models.Medicine
	if err := config.DB.First(&medicine, id).Error; err != nil {
		return err
	}
	medicine.Quantity += quantity
	if medicine.Quantity < 0 {
		medicine.Quantity = 0
	}
	return config.DB.Save(&medicine).Error
}

func (r *MedicineRepository) FindLowStock(threshold int) ([]models.Medicine, error) {
	var medicines []models.Medicine
	err := config.DB.Where("quantity <= ?", threshold).Find(&medicines).Error
	return medicines, err
}
