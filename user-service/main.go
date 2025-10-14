package main

import (
	"log"
	"user-service/connection"
	"user-service/controllers"
	"user-service/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	_ = godotenv.Load(".env")

	// Connect to DB
	db := connection.Connect()

	auth := handlers.AuthHandler{DB: db}
	r := gin.Default()

	// ✅ Setup CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000", "*"}, // frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Routes
	r.POST("/register", auth.Register)
	r.POST("/login", auth.Login)
	r.POST("/logout", handlers.AuthMiddleware(), controllers.Logout)
	sc := controllers.SubscriptionController{DB: db}
	wc := controllers.WatchlistController{DB: db}

	protected := r.Group("/")
	protected.Use(handlers.AuthMiddleware())
	{
		protected.GET("/profile", controllers.GetProfile(db))
		protected.PATCH("/profile", controllers.UpdateProfile(db))
		protected.PATCH("/subscribe", sc.UpdateUserSubscription)

		protected.PATCH("/profile/password", controllers.ChangePassword(db))

		protected.POST("/profile/watchlist", wc.AddToWatchlist)
		protected.GET("/profile/watchlist", wc.GetWatchlist)
		protected.DELETE("/profile/watchlist/:movieId", wc.RemoveFromWatchlist)
	}


	log.Println("User service running on :8001")
	r.Run(":8001")
}
