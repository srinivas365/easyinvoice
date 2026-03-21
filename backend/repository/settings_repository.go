package repository

import (
	"pharmacy-erp/config"
	"pharmacy-erp/models"
)

type SettingsRepository struct{}

func NewSettingsRepository() *SettingsRepository {
	return &SettingsRepository{}
}

func (r *SettingsRepository) Get() (*models.Settings, error) {
	var settings models.Settings
	result := config.DB.FirstOrCreate(&settings, models.Settings{ID: 1})
	return &settings, result.Error
}

func (r *SettingsRepository) Update(settings *models.Settings) error {
	settings.ID = 1 // Always update the first (and only) settings record
	return config.DB.Save(settings).Error
}
