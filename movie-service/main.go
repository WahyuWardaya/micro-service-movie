package main

import (
	"log"
	"movie-service/connection"
	"movie-service/controllers"
	"movie-service/handlers"

	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")

	db := connection.Connect()
	mc := controllers.MovieController{DB: db}
	gc := controllers.GenreController{DB: db}
	ac := controllers.ActorController{DB: db}

	r := gin.Default()

	// Movie public
	r.GET("/movies", mc.GetMovies)
	r.GET("/movies/:id", mc.GetMovieByID)

	// Genre public
	r.GET("/genres", gc.ListGenres)
	r.GET("/genres/:id", gc.GetGenre)

	// Actor public
	r.GET("/actors", ac.ListActors)
	r.GET("/actors/:id", ac.GetActor)

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000", "*"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Protected endpoints
	protected := r.Group("/")
	protected.Use(handlers.AuthMiddleware())
	{
		protected.POST("/movies", mc.CreateMovie)
		protected.PATCH("/movies/:id", mc.UpdateMovie)
		protected.DELETE("/movies/:id", mc.DeleteMovie)

		protected.POST("/genres", gc.CreateGenre)
		protected.PATCH("/genres/:id", gc.UpdateGenre)
		protected.DELETE("/genres/:id", gc.DeleteGenre)

		protected.POST("/actors", ac.CreateActor)
		protected.PATCH("/actors/:id", ac.UpdateActor)
		protected.DELETE("/actors/:id", ac.DeleteActor)
	}

	port := os.Getenv("MOVIE_SERVICE_PORT")
	if port == "" {
		port = "8002"
	}
	log.Println("Movie service running on :" + port)
	r.Run(":" + port)
}
