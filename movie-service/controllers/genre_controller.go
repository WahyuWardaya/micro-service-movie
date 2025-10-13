package controllers

import (
	"net/http"
	"strconv"

	"movie-service/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GenreController struct {
	DB *gorm.DB
}

type createGenreRequest struct {
	Name string `json:"name" binding:"required"`
}

type updateGenreRequest struct {
	Name *string `json:"name"`
}

// POST /genres  (auth required)
func (gc *GenreController) CreateGenre(c *gin.Context) {
	var req createGenreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check duplicate
	var existing models.Genre
	if err := gc.DB.Where("LOWER(name) = ?", req.Name).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "genre already exists"})
		return
	}

	genre := models.Genre{
		Name: req.Name,
	}
	if err := gc.DB.Create(&genre).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create genre"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": genre.ID, "name": genre.Name})
}

// GET /genres  (public)
func (gc *GenreController) ListGenres(c *gin.Context) {
	var genres []models.Genre
	if err := gc.DB.Find(&genres).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list genres"})
		return
	}
	c.JSON(http.StatusOK, genres)
}

// GET /genres/:id (public)
func (gc *GenreController) GetGenre(c *gin.Context) {
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var genre models.Genre
	if err := gc.DB.First(&genre, uint(id64)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "genre not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query genre"})
		return
	}
	c.JSON(http.StatusOK, genre)
}

// PATCH /genres/:id (auth required)
func (gc *GenreController) UpdateGenre(c *gin.Context) {
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req updateGenreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var genre models.Genre
	if err := gc.DB.First(&genre, uint(id64)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "genre not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query genre"})
		return
	}

	if req.Name != nil {
		genre.Name = *req.Name
	}

	if err := gc.DB.Save(&genre).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update genre"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "genre updated", "genre": genre})
}

// DELETE /genres/:id  (auth required)
func (gc *GenreController) DeleteGenre(c *gin.Context) {
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := gc.DB.Delete(&models.Genre{}, uint(id64)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete genre"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "genre deleted"})
}
