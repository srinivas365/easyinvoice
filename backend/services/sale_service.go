package services

import (
	"errors"
	"pharmacy-erp/config"
	"pharmacy-erp/models"
	"pharmacy-erp/repository"
	"time"
)

type SaleService struct {
	saleRepo     *repository.SaleRepository
	medicineRepo *repository.MedicineRepository
}

func NewSaleService() *SaleService {
	return &SaleService{
		saleRepo:     repository.NewSaleRepository(),
		medicineRepo: repository.NewMedicineRepository(),
	}
}

type CreateSaleRequest struct {
	CustomerName   string           `json:"customer_name"`
	CustomerNumber string           `json:"customer_number"`
	PrescribedBy   string           `json:"prescribed_by"`
	Items          []SaleItemRequest `json:"items" binding:"required"`
}

type SaleItemRequest struct {
	MedicineID   *uint    `json:"medicine_id"`   // Optional for manual entries
	MedicineName string   `json:"medicine_name"`  // Particulars; required for manual entries
	SNo          int      `json:"s_no"`
	Pack         string   `json:"pack"`
	Quantity     int      `json:"quantity" binding:"required,min=1"`
	Batch        string   `json:"batch"`
	Expiry       string   `json:"expiry"` // ISO date for manual entry
	MRP          float64  `json:"mrp"`
	Price        float64  `json:"price"` // Unit price; for manual entry can be derived from amount/quantity
	Discount     float64  `json:"discount"`
	DiscountType string   `json:"discount_type"`
	Total        float64  `json:"total"` // Amount (line total); if set used for manual entry
}

func (s *SaleService) CreateSale(req CreateSaleRequest) (*models.Sale, error) {
	var totalAmount float64
	var saleItems []models.SaleItem

	for _, item := range req.Items {
		var medicineID *uint
		var medicineName string
		var itemPrice float64

		// Check if it's a manual entry or from medicine catalog
		if item.MedicineID != nil && *item.MedicineID > 0 {
			// From medicine catalog - validate stock
			medicine, err := s.medicineRepo.FindByID(*item.MedicineID)
			if err != nil {
				return nil, errors.New("medicine not found")
			}

			if medicine.Quantity < item.Quantity {
				return nil, errors.New("insufficient stock for " + medicine.Name)
			}

			medicineID = item.MedicineID
			medicineName = medicine.Name
			itemPrice = item.Price // Use provided price or default to selling price
			if itemPrice == 0 {
				itemPrice = medicine.SellingPrice
			}

			// Reduce stock
			err = s.medicineRepo.UpdateStock(*item.MedicineID, -item.Quantity)
			if err != nil {
				return nil, err
			}
		} else {
			// Manual entry
			if item.MedicineName == "" {
				return nil, errors.New("medicine name (particulars) is required for manual entries")
			}
			medicineName = item.MedicineName
			if item.Total > 0 {
				itemPrice = item.Total / float64(item.Quantity)
			} else {
				itemPrice = item.Price
				if itemPrice <= 0 {
					return nil, errors.New("price or amount must be greater than 0 for manual entries")
				}
			}
		}

		// Calculate item total: use Total (amount) if provided for manual entry, else from price and discount
		var itemTotal float64
		discountType := item.DiscountType
		if discountType == "" {
			discountType = "amount"
		}
		if item.MedicineID == nil || *item.MedicineID == 0 {
			// Manual entry: Amount = Quantity × MRP
			if item.MRP > 0 {
				itemTotal = float64(item.Quantity) * item.MRP
			} else if item.Total > 0 {
				itemTotal = item.Total
			} else {
				subtotal := float64(item.Quantity) * itemPrice
				var discountAmount float64
				if discountType == "percentage" {
					if item.Discount > 100 {
						item.Discount = 100
					}
					discountAmount = subtotal * (item.Discount / 100)
				} else {
					discountAmount = item.Discount
					if discountAmount < 0 {
						discountAmount = 0
					}
				}
				itemTotal = subtotal - discountAmount
				if itemTotal < 0 {
					itemTotal = 0
				}
			}
		} else {
			subtotal := float64(item.Quantity) * itemPrice
			var discountAmount float64
			if discountType == "percentage" {
				if item.Discount > 100 {
					item.Discount = 100
				}
				discountAmount = subtotal * (item.Discount / 100)
			} else {
				discountAmount = item.Discount
				if discountAmount < 0 {
					discountAmount = 0
				}
			}
			itemTotal = subtotal - discountAmount
			if itemTotal < 0 {
				itemTotal = 0
			}
		}
		totalAmount += itemTotal

		var expiry *time.Time
		if item.Expiry != "" {
			if t, err := time.Parse("2006-01-02", item.Expiry); err == nil {
				expiry = &t
			}
		}

		saleItems = append(saleItems, models.SaleItem{
			MedicineID:   medicineID,
			MedicineName: medicineName,
			SNo:          item.SNo,
			Pack:         item.Pack,
			Quantity:     item.Quantity,
			Batch:        item.Batch,
			Expiry:       expiry,
			MRP:          item.MRP,
			Price:        itemPrice,
			Discount:     item.Discount,
			DiscountType: discountType,
			Total:        itemTotal,
		})
	}

	sale := &models.Sale{
		CustomerName:   req.CustomerName,
		CustomerNumber: req.CustomerNumber,
		PrescribedBy:   req.PrescribedBy,
		TotalAmount:    totalAmount,
		Items:          saleItems,
	}

	err := s.saleRepo.Create(sale)
	if err != nil {
		return nil, err
	}

	// Reload with items
	return s.saleRepo.FindByID(sale.ID)
}

func (s *SaleService) GetAllSales() ([]models.Sale, error) {
	return s.saleRepo.FindAll()
}

func (s *SaleService) GetArchivedSales() ([]models.Sale, error) {
	return s.saleRepo.FindAllArchived()
}

func (s *SaleService) RestoreSale(id uint) error {
	// Get archived sale using unscoped
	archivedSale, err := s.saleRepo.FindByIDUnscoped(id)
	if err != nil {
		return errors.New("archived sale not found")
	}
	
	// Check if it's actually deleted
	if archivedSale.DeletedAt.Time.IsZero() {
		return errors.New("sale is not archived")
	}
	
	// Restore the sale
	err = s.saleRepo.Restore(id)
	if err != nil {
		return err
	}
	
	// Get the restored sale to restore stock
	restoredSale, err := s.saleRepo.FindByID(id)
	if err != nil {
		return err
	}
	
	// Reduce stock from all items (only if they were from medicine catalog)
	// This reverses the stock restoration that happened during deletion
	for _, item := range restoredSale.Items {
		if item.MedicineID != nil {
			err = s.medicineRepo.UpdateStock(*item.MedicineID, -item.Quantity)
			if err != nil {
				return err
			}
		}
	}
	
	return nil
}

func (s *SaleService) GetSaleByID(id uint) (*models.Sale, error) {
	return s.saleRepo.FindByID(id)
}

func (s *SaleService) GetTodaySales() (float64, error) {
	return s.saleRepo.GetTodaySales()
}

func (s *SaleService) GetTotalRevenue() (float64, error) {
	return s.saleRepo.GetTotalRevenue()
}

func (s *SaleService) UpdateSale(id uint, req CreateSaleRequest) (*models.Sale, error) {
	// Get existing sale
	existingSale, err := s.saleRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("sale not found")
	}

	// Restore stock from existing items (only if they were from medicine catalog)
	for _, item := range existingSale.Items {
		if item.MedicineID != nil {
			err = s.medicineRepo.UpdateStock(*item.MedicineID, item.Quantity)
			if err != nil {
				return nil, err
			}
		}
	}

	// Hard delete existing items (use Unscoped to permanently delete)
	err = config.DB.Unscoped().Where("sale_id = ?", id).Delete(&models.SaleItem{}).Error
	if err != nil {
		return nil, err
	}

	// Create new items (similar to CreateSale)
	var totalAmount float64
	var saleItems []models.SaleItem

	for _, item := range req.Items {
		var medicineID *uint
		var medicineName string
		var itemPrice float64

		// Check if it's a manual entry or from medicine catalog
		if item.MedicineID != nil && *item.MedicineID > 0 {
			// From medicine catalog - validate stock
			medicine, err := s.medicineRepo.FindByID(*item.MedicineID)
			if err != nil {
				return nil, errors.New("medicine not found")
			}

			if medicine.Quantity < item.Quantity {
				return nil, errors.New("insufficient stock for " + medicine.Name)
			}

			medicineID = item.MedicineID
			medicineName = medicine.Name
			itemPrice = item.Price
			if itemPrice == 0 {
				itemPrice = medicine.SellingPrice
			}

			// Reduce stock
			err = s.medicineRepo.UpdateStock(*item.MedicineID, -item.Quantity)
			if err != nil {
				return nil, err
			}
		} else {
			// Manual entry
			if item.MedicineName == "" {
				return nil, errors.New("medicine name (particulars) is required for manual entries")
			}
			medicineName = item.MedicineName
			if item.Total > 0 {
				itemPrice = item.Total / float64(item.Quantity)
			} else {
				itemPrice = item.Price
				if itemPrice <= 0 {
					return nil, errors.New("price or amount must be greater than 0 for manual entries")
				}
			}
		}

		var itemTotal float64
		discountType := item.DiscountType
		if discountType == "" {
			discountType = "amount"
		}
		if item.MedicineID == nil || *item.MedicineID == 0 {
			// Manual entry: Amount = Quantity × MRP
			if item.MRP > 0 {
				itemTotal = float64(item.Quantity) * item.MRP
			} else if item.Total > 0 {
				itemTotal = item.Total
			} else {
				subtotal := float64(item.Quantity) * itemPrice
				var discountAmount float64
				if discountType == "percentage" {
					if item.Discount > 100 {
						item.Discount = 100
					}
					discountAmount = subtotal * (item.Discount / 100)
				} else {
					discountAmount = item.Discount
					if discountAmount < 0 {
						discountAmount = 0
					}
				}
				itemTotal = subtotal - discountAmount
				if itemTotal < 0 {
					itemTotal = 0
				}
			}
		} else {
			subtotal := float64(item.Quantity) * itemPrice
			var discountAmount float64
			if discountType == "percentage" {
				if item.Discount > 100 {
					item.Discount = 100
				}
				discountAmount = subtotal * (item.Discount / 100)
			} else {
				discountAmount = item.Discount
				if discountAmount < 0 {
					discountAmount = 0
				}
			}
			itemTotal = subtotal - discountAmount
			if itemTotal < 0 {
				itemTotal = 0
			}
		}
		totalAmount += itemTotal

		var expiry *time.Time
		if item.Expiry != "" {
			if t, err := time.Parse("2006-01-02", item.Expiry); err == nil {
				expiry = &t
			}
		}

		saleItems = append(saleItems, models.SaleItem{
			ID:           0,
			SaleID:       id,
			MedicineID:   medicineID,
			MedicineName: medicineName,
			SNo:          item.SNo,
			Pack:         item.Pack,
			Quantity:     item.Quantity,
			Batch:        item.Batch,
			Expiry:       expiry,
			MRP:          item.MRP,
			Price:        itemPrice,
			Discount:     item.Discount,
			DiscountType: discountType,
			Total:        itemTotal,
		})
	}

	// Update sale: clear Items so GORM does not re-save old items (we already deleted them and will create new ones below)
	existingSale.CustomerName = req.CustomerName
	existingSale.CustomerNumber = req.CustomerNumber
	existingSale.PrescribedBy = req.PrescribedBy
	existingSale.TotalAmount = totalAmount
	existingSale.Items = nil

	err = s.saleRepo.Update(existingSale)
	if err != nil {
		return nil, err
	}

	// Create new items (old items were already hard-deleted above)
	for i := range saleItems {
		err = config.DB.Create(&saleItems[i]).Error
		if err != nil {
			return nil, err
		}
	}

	// Reload with items
	return s.saleRepo.FindByID(id)
}

func (s *SaleService) DeleteSale(id uint) error {
	// Get existing sale
	existingSale, err := s.saleRepo.FindByID(id)
	if err != nil {
		return errors.New("sale not found")
	}

	// Restore stock from all items (only if they were from medicine catalog)
	for _, item := range existingSale.Items {
		if item.MedicineID != nil {
			err = s.medicineRepo.UpdateStock(*item.MedicineID, item.Quantity)
			if err != nil {
				return err
			}
		}
	}

	// Delete sale (items will be cascade deleted)
	return s.saleRepo.Delete(id)
}
