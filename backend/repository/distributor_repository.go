package repository

import (
	"pharmacy-erp/config"
	"pharmacy-erp/models"
)

type DistributorRepository struct{}

func NewDistributorRepository() *DistributorRepository {
	return &DistributorRepository{}
}

func (r *DistributorRepository) Create(distributor *models.Distributor) error {
	return config.DB.Create(distributor).Error
}

func (r *DistributorRepository) FindAll() ([]models.Distributor, error) {
	var distributors []models.Distributor
	err := config.DB.Find(&distributors).Error
	return distributors, err
}

func (r *DistributorRepository) FindByID(id uint) (*models.Distributor, error) {
	var distributor models.Distributor
	err := config.DB.First(&distributor, id).Error
	if err != nil {
		return nil, err
	}
	return &distributor, nil
}

func (r *DistributorRepository) Update(distributor *models.Distributor) error {
	return config.DB.Save(distributor).Error
}

func (r *DistributorRepository) Delete(id uint) error {
	return config.DB.Delete(&models.Distributor{}, id).Error
}
