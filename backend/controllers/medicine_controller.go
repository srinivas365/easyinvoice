package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"pharmacy-erp/models"
	"pharmacy-erp/services"
)

type MedicineController struct {
	medicineService *services.MedicineService
}

func NewMedicineController() *MedicineController {
	return &MedicineController{
		medicineService: services.NewMedicineService(),
	}
}

func (c *MedicineController) Create(ctx *gin.Context) {
	var medicine models.Medicine
	if err := ctx.ShouldBindJSON(&medicine); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.medicineService.CreateMedicine(&medicine); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, medicine)
}

func (c *MedicineController) GetAll(ctx *gin.Context) {
	search := ctx.Query("search")
	medicines, err := c.medicineService.GetAllMedicines(search)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, medicines)
}

func (c *MedicineController) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	medicine, err := c.medicineService.GetMedicineByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "medicine not found"})
		return
	}

	ctx.JSON(http.StatusOK, medicine)
}

func (c *MedicineController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var medicine models.Medicine
	if err := ctx.ShouldBindJSON(&medicine); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.medicineService.UpdateMedicine(uint(id), &medicine); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updated, _ := c.medicineService.GetMedicineByID(uint(id))
	ctx.JSON(http.StatusOK, updated)
}

func (c *MedicineController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := c.medicineService.DeleteMedicine(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "medicine deleted successfully"})
}
