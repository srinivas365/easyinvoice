package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"pharmacy-erp/config"
	"pharmacy-erp/models"
	"pharmacy-erp/services"
)

type DashboardController struct {
	medicineService *services.MedicineService
	saleService     *services.SaleService
}

func NewDashboardController() *DashboardController {
	return &DashboardController{
		medicineService: services.NewMedicineService(),
		saleService:     services.NewSaleService(),
	}
}

func (c *DashboardController) GetStats(ctx *gin.Context) {
	// Total medicines
	var totalMedicines int64
	config.DB.Model(&models.Medicine{}).Count(&totalMedicines)

	// Low stock medicines (threshold: 10)
	lowStockMedicines, _ := c.medicineService.GetLowStockMedicines(10)

	// Today sales
	todaySales, _ := c.saleService.GetTodaySales()

	// Total revenue
	totalRevenue, _ := c.saleService.GetTotalRevenue()

	ctx.JSON(http.StatusOK, gin.H{
		"total_medicines":     totalMedicines,
		"low_stock_medicines": len(lowStockMedicines),
		"low_stock_list":      lowStockMedicines,
		"today_sales":         todaySales,
		"total_revenue":       totalRevenue,
	})
}
