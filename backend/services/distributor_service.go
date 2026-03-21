package services

import (
	"errors"
	"pharmacy-erp/models"
	"pharmacy-erp/repository"
)

type DistributorService struct {
	distributorRepo *repository.DistributorRepository
}

func NewDistributorService() *DistributorService {
	return &DistributorService{
		distributorRepo: repository.NewDistributorRepository(),
	}
}

func (s *DistributorService) CreateDistributor(distributor *models.Distributor) error {
	return s.distributorRepo.Create(distributor)
}

func (s *DistributorService) GetAllDistributors() ([]models.Distributor, error) {
	return s.distributorRepo.FindAll()
}

func (s *DistributorService) GetDistributorByID(id uint) (*models.Distributor, error) {
	return s.distributorRepo.FindByID(id)
}

func (s *DistributorService) UpdateDistributor(id uint, distributor *models.Distributor) error {
	existing, err := s.distributorRepo.FindByID(id)
	if err != nil {
		return errors.New("distributor not found")
	}

	existing.Name = distributor.Name
	existing.ContactPerson = distributor.ContactPerson
	existing.Phone = distributor.Phone
	existing.Email = distributor.Email
	existing.Address = distributor.Address

	return s.distributorRepo.Update(existing)
}

func (s *DistributorService) DeleteDistributor(id uint) error {
	return s.distributorRepo.Delete(id)
}
