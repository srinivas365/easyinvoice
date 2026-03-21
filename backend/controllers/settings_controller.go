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
	settings, err := c.settingsService.GetSettings()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, settings)
}

func (c *SettingsController) Update(ctx *gin.Context) {
	var settings models.Settings
	if err := ctx.ShouldBindJSON(&settings); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.settingsService.UpdateSettings(&settings); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, settings)
}
