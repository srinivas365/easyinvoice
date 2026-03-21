package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"pharmacy-erp/models"
	"pharmacy-erp/services"
)

type DistributorController struct {
	distributorService *services.DistributorService
}

func NewDistributorController() *DistributorController {
	return &DistributorController{
		distributorService: services.NewDistributorService(),
	}
}

func (c *DistributorController) Create(ctx *gin.Context) {
	var distributor models.Distributor
	if err := ctx.ShouldBindJSON(&distributor); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.distributorService.CreateDistributor(&distributor); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, distributor)
}

func (c *DistributorController) GetAll(ctx *gin.Context) {
	distributors, err := c.distributorService.GetAllDistributors()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, distributors)
}

func (c *DistributorController) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	distributor, err := c.distributorService.GetDistributorByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "distributor not found"})
		return
	}

	ctx.JSON(http.StatusOK, distributor)
}

func (c *DistributorController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var distributor models.Distributor
	if err := ctx.ShouldBindJSON(&distributor); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.distributorService.UpdateDistributor(uint(id), &distributor); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updated, _ := c.distributorService.GetDistributorByID(uint(id))
	ctx.JSON(http.StatusOK, updated)
}

func (c *DistributorController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := c.distributorService.DeleteDistributor(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "distributor deleted successfully"})
}
