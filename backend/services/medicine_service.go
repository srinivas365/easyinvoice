package services

import (
	"errors"
	"pharmacy-erp/models"
	"pharmacy-erp/repository"
)

type MedicineService struct {
	medicineRepo *repository.MedicineRepository
}

func NewMedicineService() *MedicineService {
	return &MedicineService{
		medicineRepo: repository.NewMedicineRepository(),
	}
}

func (s *MedicineService) CreateMedicine(medicine *models.Medicine) error {
	return s.medicineRepo.Create(medicine)
}

func (s *MedicineService) GetAllMedicines(search string) ([]models.Medicine, error) {
	return s.medicineRepo.FindAll(search)
}

func (s *MedicineService) GetMedicineByID(id uint) (*models.Medicine, error) {
	return s.medicineRepo.FindByID(id)
}

func (s *MedicineService) UpdateMedicine(id uint, medicine *models.Medicine) error {
	existing, err := s.medicineRepo.FindByID(id)
	if err != nil {
		return errors.New("medicine not found")
	}

	existing.Name = medicine.Name
	existing.Manufacturer = medicine.Manufacturer
	existing.BatchNumber = medicine.BatchNumber
	existing.ExpiryDate = medicine.ExpiryDate
	existing.PurchasePrice = medicine.PurchasePrice
	existing.SellingPrice = medicine.SellingPrice
	existing.Quantity = medicine.Quantity

	return s.medicineRepo.Update(existing)
}

func (s *MedicineService) DeleteMedicine(id uint) error {
	return s.medicineRepo.Delete(id)
}

func (s *MedicineService) GetLowStockMedicines(threshold int) ([]models.Medicine, error) {
	return s.medicineRepo.FindLowStock(threshold)
}
