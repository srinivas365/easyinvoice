package services

import (
	"errors"
	"pharmacy-erp/config"
	"pharmacy-erp/models"
	"pharmacy-erp/repository"
)

type PurchaseInvoiceService struct {
	repo            *repository.PurchaseInvoiceRepository
	distributorRepo *repository.DistributorRepository
	productRepo     *repository.ProductRepository
}

func NewPurchaseInvoiceService() *PurchaseInvoiceService {
	return &PurchaseInvoiceService{
		repo:            repository.NewPurchaseInvoiceRepository(),
		distributorRepo: repository.NewDistributorRepository(),
		productRepo:     repository.NewProductRepository(),
	}
}

type PurchaseInvoiceItemRequest struct {
	ProductID       uint    `json:"product_id" binding:"required"`
	BatchNo         string  `json:"batch_no"`
	Expiry          FlexDate `json:"expiry"`
	Quantity        float64 `json:"quantity" binding:"required,min=0"`
	Free            float64 `json:"free"`
	Mrp             float64 `json:"mrp"`
	Rate            float64 `json:"rate" binding:"required"`
	Amount          float64 `json:"amount" binding:"required"`
	Gst             float64 `json:"gst"`
	DiscountPercent float64 `json:"discount_percent"`
}

type CreatePurchaseInvoiceRequest struct {
	DistributorID      uint                         `json:"distributor_id" binding:"required"`
	PurchaseDate       FlexDate                     `json:"purchase_date" binding:"required"`
	AdditionalDiscount float64                      `json:"additional_discount"`
	Items              []PurchaseInvoiceItemRequest `json:"items" binding:"required,min=1"`
}

func (s *PurchaseInvoiceService) Create(req CreatePurchaseInvoiceRequest) (*models.PurchaseInvoice, error) {
	_, err := s.distributorRepo.FindByID(req.DistributorID)
	if err != nil {
		return nil, errors.New("distributor not found")
	}

	var itemsTotal float64
	var items []models.PurchaseInvoiceItem
	for _, it := range req.Items {
		if _, err := s.productRepo.FindByID(it.ProductID); err != nil {
			return nil, errors.New("product not found")
		}
		itemsTotal += it.Amount
		items = append(items, models.PurchaseInvoiceItem{
			ProductID:       it.ProductID,
			BatchNo:         it.BatchNo,
			Expiry:          it.Expiry.Time,
			Quantity:        it.Quantity,
			Free:            it.Free,
			Mrp:             it.Mrp,
			Rate:            it.Rate,
			Amount:          it.Amount,
			Gst:             it.Gst,
			DiscountPercent: it.DiscountPercent,
		})
	}

	additionalDiscount := req.AdditionalDiscount
	if additionalDiscount < 0 {
		additionalDiscount = 0
	}
	totalAmount := itemsTotal - additionalDiscount
	if totalAmount < 0 {
		totalAmount = 0
	}

	inv := &models.PurchaseInvoice{
		DistributorID:      req.DistributorID,
		PurchaseDate:       req.PurchaseDate.Time,
		AdditionalDiscount: additionalDiscount,
		TotalAmount:        totalAmount,
		Items:              items,
	}
	if err := s.repo.Create(inv); err != nil {
		return nil, err
	}
	return s.repo.FindByID(inv.ID)
}

func (s *PurchaseInvoiceService) GetAll() ([]models.PurchaseInvoice, error) {
	return s.repo.FindAll()
}

func (s *PurchaseInvoiceService) GetByID(id uint) (*models.PurchaseInvoice, error) {
	return s.repo.FindByID(id)
}

func (s *PurchaseInvoiceService) Update(id uint, req CreatePurchaseInvoiceRequest) (*models.PurchaseInvoice, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("purchase invoice not found")
	}
	_, err = s.distributorRepo.FindByID(req.DistributorID)
	if err != nil {
		return nil, errors.New("distributor not found")
	}

	var itemsTotal float64
	var items []models.PurchaseInvoiceItem
	for _, it := range req.Items {
		if _, err := s.productRepo.FindByID(it.ProductID); err != nil {
			return nil, errors.New("product not found")
		}
		itemsTotal += it.Amount
		items = append(items, models.PurchaseInvoiceItem{
			PurchaseInvoiceID: id,
			ProductID:         it.ProductID,
			BatchNo:           it.BatchNo,
			Expiry:            it.Expiry.Time,
			Quantity:          it.Quantity,
			Free:              it.Free,
			Mrp:               it.Mrp,
			Rate:              it.Rate,
			Amount:            it.Amount,
			Gst:               it.Gst,
			DiscountPercent:   it.DiscountPercent,
		})
	}

	additionalDiscount := req.AdditionalDiscount
	if additionalDiscount < 0 {
		additionalDiscount = 0
	}
	totalAmount := itemsTotal - additionalDiscount
	if totalAmount < 0 {
		totalAmount = 0
	}

	if err := s.repo.DeleteItemsByInvoiceID(id); err != nil {
		return nil, err
	}
	for i := range items {
		if err := config.DB.Create(&items[i]).Error; err != nil {
			return nil, err
		}
	}

	existing.DistributorID = req.DistributorID
	existing.PurchaseDate = req.PurchaseDate.Time
	existing.AdditionalDiscount = additionalDiscount
	existing.TotalAmount = totalAmount
	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}
	return s.repo.FindByID(id)
}

func (s *PurchaseInvoiceService) Delete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("purchase invoice not found")
	}
	return s.repo.Delete(id)
}

