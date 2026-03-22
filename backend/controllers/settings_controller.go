package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"pharmacy-erp/models"
	"pharmacy-erp/services"
)

type SettingsController struct {
	settingsService *services.SettingsService
}

func NewSettingsController() *SettingsController {
	return &SettingsController{
		settingsService: services.NewSettingsService(),
	}
}

func (c *SettingsController) Get(ctx *gin.Context) {
	accountID, ok := ctx.Get("account_id")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "account context missing"})
		return
	}
	settings, err := c.settingsService.GetSettings(accountID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, settings)
}

func (c *SettingsController) Update(ctx *gin.Context) {
	accountID, ok := ctx.Get("account_id")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "account context missing"})
		return
	}
	var settings models.Settings
	if err := ctx.ShouldBindJSON(&settings); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.settingsService.UpdateSettings(accountID.(uint), &settings); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, settings)
}
