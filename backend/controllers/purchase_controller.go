package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"pharmacy-erp/models"
	"pharmacy-erp/services"
)

type PurchaseController struct {
	purchaseService *services.PurchaseService
}

func NewPurchaseController() *PurchaseController {
	return &PurchaseController{
		purchaseService: services.NewPurchaseService(),
	}
}

func (c *PurchaseController) Create(ctx *gin.Context) {
	var purchase models.Purchase
	if err := ctx.ShouldBindJSON(&purchase); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.purchaseService.CreatePurchase(&purchase); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, purchase)
}

func (c *PurchaseController) GetAll(ctx *gin.Context) {
	purchases, err := c.purchaseService.GetAllPurchases()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, purchases)
}

func (c *PurchaseController) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	purchase, err := c.purchaseService.GetPurchaseByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "purchase not found"})
		return
	}

	ctx.JSON(http.StatusOK, purchase)
}
