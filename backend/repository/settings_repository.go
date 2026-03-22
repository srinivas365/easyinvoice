package repository

import (
	"errors"
	"strings"

	"pharmacy-erp/config"
	"pharmacy-erp/models"

	"gorm.io/gorm"
)

type SettingsRepository struct{}

func NewSettingsRepository() *SettingsRepository {
	return &SettingsRepository{}
}

func isSettingsAccountUniqueErr(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique") && strings.Contains(msg, "account_id")
}

func (r *SettingsRepository) GetForAccount(accountID uint) (*models.Settings, error) {
	var settings models.Settings
	if err := config.DB.Where("account_id = ?", accountID).First(&settings).Error; err == nil {
		return &settings, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Legacy row: account_id unset (0 / NULL) does not match the query above — adopt it instead of inserting.
	var orphan models.Settings
	err := config.DB.Where("account_id IS NULL OR account_id = ?", 0).Order("id").Limit(1).First(&orphan).Error
	if err == nil {
		orphan.AccountID = accountID
		if orphan.CompanyName == "" {
			orphan.CompanyName = "My Pharmacy"
		}
		if orphan.LicenseID == "" {
			orphan.LicenseID = "N/A"
		}
		if err := config.DB.Save(&orphan).Error; err != nil {
			if isSettingsAccountUniqueErr(err) {
				if err2 := config.DB.Where("account_id = ?", accountID).First(&settings).Error; err2 == nil {
					return &settings, nil
				}
			}
			return nil, err
		}
		return &orphan, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	settings = models.Settings{
		AccountID:   accountID,
		CompanyName: "My Pharmacy",
		LicenseID:   "N/A",
	}
	if err := config.DB.Create(&settings).Error; err != nil {
		if isSettingsAccountUniqueErr(err) {
			if err2 := config.DB.Where("account_id = ?", accountID).First(&settings).Error; err2 == nil {
				return &settings, nil
			}
		}
		return nil, err
	}
	return &settings, nil
}

func (r *SettingsRepository) UpdateForAccount(accountID uint, settings *models.Settings) error {
	existing, err := r.GetForAccount(accountID)
	if err != nil {
		return err
	}
	settings.ID = existing.ID
	settings.AccountID = accountID
	return config.DB.Save(settings).Error
}
