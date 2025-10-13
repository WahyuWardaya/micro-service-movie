package main

import (
	"log"
	"os"
	"subscription-service/connection"
	"subscription-service/controllers"
	"subscription-service/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")

	db := connection.Connect()
	sc := controllers.SubscriptionController{DB: db}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000", "*"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Protected routes
	protected := r.Group("/")
	protected.Use(handlers.AuthMiddleware())
	{
		protected.POST("/subscribe", sc.CreateSubscription)       // buy subscription
		protected.GET("/subscriptions/me", sc.GetMySubscriptions) // list subscriptions for current user (optional, depends on your controller)
	}

	port := os.Getenv("SUBSCRIPTION_SERVICE_PORT")
	if port == "" {
		port = "8003"
	}
	log.Println("Subscription service running on :" + port)
	r.Run(":" + port)
}
