package repository

import (
	"pharmacy-erp/config"
	"pharmacy-erp/models"
)

type ProductRepository struct{}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{}
}

func (r *ProductRepository) Create(p *models.Product) error {
	return config.DB.Create(p).Error
}

func (r *ProductRepository) FindAll(search string) ([]models.Product, error) {
	var list []models.Product
	q := config.DB.Order("name")
	if search != "" {
		q = q.Where("name LIKE ? OR mfg LIKE ? OR hsn_code LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	err := q.Find(&list).Error
	return list, err
}

func (r *ProductRepository) FindByID(id uint) (*models.Product, error) {
	var p models.Product
	err := config.DB.First(&p, id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) FindOrCreateByNameMfgHsnPack(name, mfg, hsn, pack string) (*models.Product, error) {
	var p models.Product
	err := config.DB.Where("name = ? AND mfg = ? AND hsn_code = ? AND pack = ?", name, mfg, hsn, pack).First(&p).Error
	if err == nil {
		return &p, nil
	}
	p = models.Product{Name: name, Mfg: mfg, HsnCode: hsn, Pack: pack}
	if err := config.DB.Create(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Update(p *models.Product) error {
	return config.DB.Save(p).Error
}

func (r *ProductRepository) Delete(id uint) error {
	return config.DB.Delete(&models.Product{}, id).Error
}

// ProductWithQuantity is used for Medicines list (product + aggregated quantity from purchase invoice items)
type ProductWithQuantity struct {
	ID        uint    `json:"id"`
	Name      string  `json:"name"`
	Mfg       string  `json:"mfg"`
	HsnCode   string  `json:"hsn_code"`
	Pack      string  `json:"pack"`
	Quantity  float64 `json:"quantity"`
}

func (r *ProductRepository) FindAllWithQuantity(search string) ([]ProductWithQuantity, error) {
	query := `
		SELECT p.id, p.name, p.mfg, p.hsn_code, p.pack,
		       COALESCE(SUM(i.quantity + i.free), 0) as quantity
		FROM products p
		LEFT JOIN purchase_invoice_items i ON i.product_id = p.id AND i.deleted_at IS NULL
		  AND i.purchase_invoice_id IN (SELECT id FROM purchase_invoices WHERE deleted_at IS NULL)
		WHERE p.deleted_at IS NULL
	`
	args := []interface{}{}
	if search != "" {
		query += ` AND (p.name LIKE ? OR p.mfg LIKE ? OR p.hsn_code LIKE ?)`
		args = append(args, "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	query += ` GROUP BY p.id ORDER BY p.name`
	results := []ProductWithQuantity{}
	err := config.DB.Raw(query, args...).Scan(&results).Error
	if err != nil {
		return nil, err
	}
	if results == nil {
		results = []ProductWithQuantity{}
	}
	return results, nil
}
