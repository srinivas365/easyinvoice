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

func (s *SettingsService) GetSettings() (*models.Settings, error) {
	return s.settingsRepo.Get()
}

func (s *SettingsService) UpdateSettings(settings *models.Settings) error {
	return s.settingsRepo.Update(settings)
}
