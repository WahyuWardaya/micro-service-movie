package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"subscription-service/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SubscriptionController struct {
    DB *gorm.DB
}

type createSubReq struct {
    Plan string `json:"plan" binding:"required"` // "monthly", "3months", "yearly"
}

// price map (in smallest currency unit, e.g. cents or in rupiah)
var priceMap = map[string]int64{
    "monthly": 45000,
    "3months": 125000,
    "yearly": 1620000,
}

// POST /subscribe (auth required — user JWT forwarded from client)
func (sc *SubscriptionController) CreateSubscription(c *gin.Context) {
    // get user id from token via middleware
    uidv, _ := c.Get("user_id")
    userID := uidv.(uint)

    var req createSubReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    price, ok := priceMap[req.Plan]
    if !ok {
        c.JSON(http.StatusBadRequest, gin.H{"error":"unknown plan"})
        return
    }

    // compute start and end
    start := time.Now()
    var end time.Time
    switch req.Plan {
    case "monthly":
        end = start.AddDate(0,1,0)
    case "3months":
        end = start.AddDate(0,3,0)
    case "yearly":
        end = start.AddDate(1,0,0)
    }

    sub := models.Subscription{
        UserID: userID,
        Plan: req.Plan,
        Amount: price,
        Status: "success",
        StartedAt: start,
        EndAt: end,
    }
    if err := sc.DB.Create(&sub).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error":"failed to create subscription"})
        return
    }

    // Now update user-service subscription by calling its endpoint and forwarding user's JWT
    userSvc := os.Getenv("USER_SERVICE_URL") // e.g., http://user-service:8001
    upd := map[string]string{
        "subscription_type": req.Plan,
        "subscription_expired_at": end.Format(time.RFC3339),
    }
    b, _ := json.Marshal(upd)
    client := &http.Client{ Timeout: 5 * time.Second }

    // forward the same Authorization header the client used
    authHeader := c.GetHeader("Authorization")
    req2, _ := http.NewRequest("PATCH", userSvc+"/subscribe", bytes.NewReader(b))
    req2.Header.Set("Content-Type","application/json")
    if authHeader != "" {
        req2.Header.Set("Authorization", authHeader)
    }

    resp, err := client.Do(req2)
    if err != nil || resp.StatusCode >= 400 {
        // warn but keep subscription record (depending on policy you might rollback)
        c.JSON(http.StatusOK, gin.H{
            "message":"subscription created but failed to update user-service",
            "subscription": sub,
        })
        return
    }

    // success
    c.JSON(http.StatusCreated, gin.H{
        "message":"subscription successful",
        "subscription": sub,
    })
}

func (sc *SubscriptionController) GetMySubscriptions(c *gin.Context) {
    uidv, _ := c.Get("user_id")
    userID := uidv.(uint)

    var subs []models.Subscription
    if err := sc.DB.Where("user_id = ?", userID).Order("start_at DESC").Find(&subs).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error":"failed to fetch subscriptions"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "user_id": userID,
        "subscriptions": subs,
    })
}