package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"pharmacy-erp/services"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{
		authService: services.NewAuthService(),
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := c.authService.Login(req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

func (c *AuthController) GetProfile(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	accountID, _ := ctx.Get("account_id")
	username, _ := ctx.Get("username")
	role, _ := ctx.Get("role")

	ctx.JSON(http.StatusOK, gin.H{
		"id":          userID,
		"account_id":  accountID,
		"username":    username,
		"role":        role,
	})
}
