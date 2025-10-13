package handlers

import (
	"net/http"
	"strings"
	"user-service/models"
	"user-service/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB *gorm.DB
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing models.User
	if err := h.DB.Where("email = ?", strings.ToLower(req.Email)).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
		return
	}

	hash, _ := utils.HashPassword(req.Password)

	user := models.User{
		Name:         req.Name,
		Email:        strings.ToLower(req.Email),
		PasswordHash: hash,
	}

	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": user.ID, "name": user.Name, "email": user.Email})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.DB.Where("email = ?", strings.ToLower(req.Email)).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, _ := utils.GenerateToken(int(user.ID), user.Email)
	c.JSON(http.StatusOK, gin.H{"token": token, "name": user.Name, "email": user.Email})
}
