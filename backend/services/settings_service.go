package services

import (
	"pharmacy-erp/models"
	"pharmacy-erp/repository"
)

type SettingsService struct {
	settingsRepo *repository.SettingsRepository
}

func NewSettingsService() *SettingsService {
	return &SettingsService{
		settingsRepo: repository.NewSettingsRepository(),
	}
}

func (s *SettingsService) GetSettings(accountID uint) (*models.Settings, error) {
	return s.settingsRepo.GetForAccount(accountID)
}

func (s *SettingsService) UpdateSettings(accountID uint, settings *models.Settings) error {
	return s.settingsRepo.UpdateForAccount(accountID, settings)
}
