package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"pharmacy-erp/services"
)

type PurchaseInvoiceController struct {
	service *services.PurchaseInvoiceService
}

func NewPurchaseInvoiceController() *PurchaseInvoiceController {
	return &PurchaseInvoiceController{
		service: services.NewPurchaseInvoiceService(),
	}
}

func (c *PurchaseInvoiceController) Create(ctx *gin.Context) {
	var req services.CreatePurchaseInvoiceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	inv, err := c.service.Create(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, inv)
}

func (c *PurchaseInvoiceController) GetAll(ctx *gin.Context) {
	list, err := c.service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, list)
}

func (c *PurchaseInvoiceController) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	inv, err := c.service.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "purchase invoice not found"})
		return
	}
	ctx.JSON(http.StatusOK, inv)
}

func (c *PurchaseInvoiceController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req services.CreatePurchaseInvoiceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	inv, err := c.service.Update(uint(id), req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, inv)
}

func (c *PurchaseInvoiceController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := c.service.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "purchase invoice deleted successfully"})
}
