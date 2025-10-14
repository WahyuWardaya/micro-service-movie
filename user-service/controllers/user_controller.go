package controllers

import (
	"net/http"
	"user-service/models"
	"user-service/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetProfile - GET /profile (auth required)
// Mengembalikan detail lengkap dari pengguna yang sedang login
func GetProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		var user models.User
		// Menggunakan First untuk mencari user berdasarkan Primary Key (ID)
		if err := db.First(&user, userID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve user"})
			return
		}

		// Mengembalikan data user tanpa password hash
		c.JSON(http.StatusOK, gin.H{
			"id":                      user.ID,
			"name":                    user.Name,
			"email":                   user.Email,
			"subscription_type":       user.SubscriptionType,
			"subscription_expired_at": user.SubscriptionExpiredAt,
		})
	}
}

type changePasswordRequest struct {
    OldPassword string `json:"old_password" binding:"required"`
    NewPassword string `json:"new_password" binding:"required"`
}

// ChangePassword - PATCH /profile/password (auth required)
func ChangePassword(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetUint("user_id")

        // 1. Bind request body ke struct
        var req changePasswordRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // 2. Ambil data user dari database
        var user models.User
        if err := db.First(&user, userID).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
            return
        }

        // 3. Verifikasi password lama (PENTING UNTUK KEAMANAN)
        if !utils.CheckPasswordHash(req.OldPassword, user.PasswordHash) {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid old password"})
            return
        }

        // 4. Hash password baru
        newHashedPassword, err := utils.HashPassword(req.NewPassword)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash new password"})
            return
        }

        // 5. Update password di database
        if err := db.Model(&user).Update("password_hash", newHashedPassword).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "password updated successfully"})
    }
}