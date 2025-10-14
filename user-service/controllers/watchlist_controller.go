package controllers

import (
	"net/http"
	"strconv"
	"user-service/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WatchlistController struct {
    DB *gorm.DB
}

// POST /profile/watchlist
func (wc *WatchlistController) AddToWatchlist(c *gin.Context) {
    userID := c.GetUint("user_id")

    var req struct {
        MovieID uint `json:"movie_id" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request, movie_id is required"})
        return
    }

    // Buat entri film jika belum ada (hanya ID-nya saja)
    movie := models.Movie{ID: req.MovieID}
    wc.DB.FirstOrCreate(&movie, movie)

    // Tambahkan asosiasi ke watchlist
    var user models.User
    if err := wc.DB.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

    if err := wc.DB.Model(&user).Association("Watchlist").Append(&movie); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add movie to watchlist"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "movie added to watchlist"})
}

// GET /profile/watchlist
func (wc *WatchlistController) GetWatchlist(c *gin.Context) {
    userID := c.GetUint("user_id")

    var user models.User
    // Menggunakan Preload untuk mengambil data film (hanya ID) dari tabel asosiasi
    if err := wc.DB.Preload("Watchlist").First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

    // Mengekstrak hanya ID film dari hasil preload
    movieIDs := make([]uint, len(user.Watchlist))
    for i, movie := range user.Watchlist {
        movieIDs[i] = movie.ID
    }
    
    // Di aplikasi nyata, Anda akan memanggil movie-service di sini untuk mendapatkan detail film
    // Untuk saat ini, kita hanya kembalikan daftar ID-nya.
    c.JSON(http.StatusOK, gin.H{
        "user_id": userID,
        "watchlist_movie_ids": movieIDs,
    })
}

// DELETE /profile/watchlist/:movieId
func (wc *WatchlistController) RemoveFromWatchlist(c *gin.Context) {
    userID := c.GetUint("user_id")
    movieIDStr := c.Param("movieId")
    
    movieID, err := strconv.ParseUint(movieIDStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
        return
    }

    var user models.User
    if err := wc.DB.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

    movie := models.Movie{ID: uint(movieID)}

    // Menghapus asosiasi dari watchlist
    if err := wc.DB.Model(&user).Association("Watchlist").Delete(&movie); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove movie from watchlist"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "movie removed from watchlist"})
}