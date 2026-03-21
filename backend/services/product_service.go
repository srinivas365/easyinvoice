package services

import (
	"pharmacy-erp/models"
	"pharmacy-erp/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService() *ProductService {
	return &ProductService{repo: repository.NewProductRepository()}
}

func (s *ProductService) Create(p *models.Product) error {
	return s.repo.Create(p)
}

func (s *ProductService) GetAll(search string) ([]models.Product, error) {
	return s.repo.FindAll(search)
}

func (s *ProductService) GetAllWithQuantity(search string) ([]repository.ProductWithQuantity, error) {
	return s.repo.FindAllWithQuantity(search)
}

func (s *ProductService) GetByID(id uint) (*models.Product, error) {
	return s.repo.FindByID(id)
}

func (s *ProductService) Update(p *models.Product) error {
	return s.repo.Update(p)
}

func (s *ProductService) Delete(id uint) error {
	return s.repo.Delete(id)
}
