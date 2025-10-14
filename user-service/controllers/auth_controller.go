package controllers

import (
	"net/http"
	"user-service/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully. Please remove token on client-side.",
	})
}

func UpdateProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id") // didapat dari middleware

		var req struct {
			Name  *string `json:"name"`
			Email *string `json:"email"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}

		// Ambil user dari DB menggunakan GORM
		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		// Update jika ada perubahan
		updates := map[string]interface{}{}
		if req.Name != nil {
			updates["name"] = *req.Name
		}
		if req.Email != nil {
			updates["email"] = *req.Email
		}

		// Jika ada field yang diupdate
		if len(updates) > 0 {
			if err := db.Model(&user).Updates(updates).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update"})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Profile updated",
			"user":    user,
		})
	}
}