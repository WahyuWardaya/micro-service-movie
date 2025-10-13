package controllers

import (
	"net/http"
	"time"
	"user-service/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PATCH /subscribe  -> protected by JWT middleware
type SubscriptionUpdateRequest struct {
    SubscriptionType string    `json:"subscription_type" binding:"required"` // monthly, 3months, yearly, none
    ExpiresAt        *time.Time `json:"expires_at"` // optional: if provided, set explicitly
}

type SubscriptionController struct {
    DB *gorm.DB
}

func (sc *SubscriptionController) UpdateUserSubscription(c *gin.Context) {
    uid := c.GetUint("user_id")
    if uid == 0 {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    var req SubscriptionUpdateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user models.User
    if err := sc.DB.First(&user, uid).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

    user.SubscriptionType = req.SubscriptionType
    if req.ExpiresAt != nil {
        user.SubscriptionExpiredAt = req.ExpiresAt
    }
    if req.SubscriptionType == "none" {
        user.SubscriptionExpiredAt = nil
    }

    if err := sc.DB.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update subscription"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "subscription updated", "subscription_type": user.SubscriptionType, "expires_at": user.SubscriptionExpiredAt})
}
