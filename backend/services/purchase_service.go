package services

import (
	"errors"
	"pharmacy-erp/models"
	"pharmacy-erp/repository"
)

type PurchaseService struct {
	purchaseRepo *repository.PurchaseRepository
	medicineRepo *repository.MedicineRepository
}

func NewPurchaseService() *PurchaseService {
	return &PurchaseService{
		purchaseRepo: repository.NewPurchaseRepository(),
		medicineRepo: repository.NewMedicineRepository(),
	}
}

func (s *PurchaseService) CreatePurchase(purchase *models.Purchase) error {
	// Verify medicine exists
	_, err := s.medicineRepo.FindByID(purchase.MedicineID)
	if err != nil {
		return errors.New("medicine not found")
	}

	// Create purchase record
	err = s.purchaseRepo.Create(purchase)
	if err != nil {
		return err
	}

	// Update stock
	return s.medicineRepo.UpdateStock(purchase.MedicineID, purchase.Quantity)
}

func (s *PurchaseService) GetAllPurchases() ([]models.Purchase, error) {
	return s.purchaseRepo.FindAll()
}

func (s *PurchaseService) GetPurchaseByID(id uint) (*models.Purchase, error) {
	return s.purchaseRepo.FindByID(id)
}
