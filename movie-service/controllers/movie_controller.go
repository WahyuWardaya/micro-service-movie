package controllers

import (
	"movie-service/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Request structs
type createMovieRequest struct {
	Title           string  `json:"title" binding:"required"`
	PosterBase64    *string `json:"poster_base64"` // optional
	DurationMinutes *int    `json:"duration_minutes"`
	Synopsis        *string `json:"synopsis"`
	ReleaseYear     *int    `json:"release_year"`
	Rating          *float32 `json:"rating"`
	Views           *int64  `json:"views"`
	Genres          []uint  `json:"genres"` // list of genre IDs
	Actors          []uint  `json:"actors"` // list of actor IDs
}

type updateMovieRequest struct {
	Title           *string  `json:"title"`
	PosterBase64    *string  `json:"poster_base64"` // if nil => keep old
	DurationMinutes *int     `json:"duration_minutes"`
	Synopsis        *string  `json:"synopsis"`
	ReleaseYear     *int     `json:"release_year"`
	Rating          *float32 `json:"rating"`
	Views           *int64   `json:"views"`
	Genres          []uint   `json:"genres"` // full replace if provided (len>0)
	Actors          []uint   `json:"actors"` // full replace if provided (len>0)
}

type MovieController struct {
	DB *gorm.DB
}

// helper to map []models.Genre -> []uint
func genreIDs(genres []models.Genre) []uint {
	out := make([]uint, 0, len(genres))
	for _, g := range genres {
		out = append(out, g.ID)
	}
	return out
}

func actorIDs(actors []models.Actor) []uint {
	out := make([]uint, 0, len(actors))
	for _, a := range actors {
		out = append(out, a.ID)
	}
	return out
}

// CreateMovie - POST /movies (auth required)
func (mc *MovieController) CreateMovie(c *gin.Context) {
	var req createMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// validate genres exist
	var genres []models.Genre
	if len(req.Genres) > 0 {
		if err := mc.DB.Find(&genres, req.Genres).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query genres"})
			return
		}
		if len(genres) != len(req.Genres) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "one or more genres not found"})
			return
		}
	}

	// validate actors exist
	var actors []models.Actor
	if len(req.Actors) > 0 {
		if err := mc.DB.Find(&actors, req.Actors).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query actors"})
			return
		}
		if len(actors) != len(req.Actors) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "one or more actors not found"})
			return
		}
	}

	movie := models.Movie{
		Title: req.Title,
	}
	if req.PosterBase64 != nil {
		movie.PosterBase64 = *req.PosterBase64
	}
	if req.DurationMinutes != nil {
		movie.DurationMinutes = *req.DurationMinutes
	}
	if req.Synopsis != nil {
		movie.Synopsis = *req.Synopsis
	}
	if req.ReleaseYear != nil {
		movie.ReleaseYear = *req.ReleaseYear
	}
	if req.Rating != nil {
		movie.Rating = *req.Rating
	}
	if req.Views != nil {
		movie.Views = *req.Views
	}

	if err := mc.DB.Create(&movie).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create movie: " + err.Error()})
		return
	}

	// set relations if any
	if len(genres) > 0 {
		if err := mc.DB.Model(&movie).Association("Genres").Replace(&genres); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to attach genres"})
			return
		}
	}
	if len(actors) > 0 {
		if err := mc.DB.Model(&movie).Association("Actors").Replace(&actors); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to attach actors"})
			return
		}
	}

	// respond with IDs for genres & actors
	var outGenres []models.Genre
	var outActors []models.Actor
	mc.DB.Model(&movie).Association("Genres").Find(&outGenres)
	mc.DB.Model(&movie).Association("Actors").Find(&outActors)

	c.JSON(http.StatusCreated, gin.H{
		"id":               movie.ID,
		"title":            movie.Title,
		"poster_base64":    movie.PosterBase64,
		"duration_minutes": movie.DurationMinutes,
		"synopsis":         movie.Synopsis,
		"release_year":     movie.ReleaseYear,
		"rating":           movie.Rating,
		"views":            movie.Views,
		"genres":           genreIDs(outGenres),
		"actors":           actorIDs(outActors),
	})
}

// GetMovies - GET /movies (public)
func (mc *MovieController) GetMovies(c *gin.Context) {
	var movies []models.Movie
	if err := mc.DB.Preload("Genres").Preload("Actors").Find(&movies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list movies"})
		return
	}

	out := make([]gin.H, 0, len(movies))
	for _, m := range movies {
		out = append(out, gin.H{
			"id":               m.ID,
			"title":            m.Title,
			"poster_base64":    m.PosterBase64,
			"duration_minutes": m.DurationMinutes,
			"synopsis":         m.Synopsis,
			"release_year":     m.ReleaseYear,
			"rating":           m.Rating,
			"views":            m.Views,
			"genres":           genreIDs(m.Genres),
			"actors":           actorIDs(m.Actors),
		})
	}
	c.JSON(http.StatusOK, out)
}

// GetMovieByID - GET /movies/:id (public)
func (mc *MovieController) GetMovieByID(c *gin.Context) {
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	id := uint(id64)

	var movie models.Movie
	if err := mc.DB.Preload("Genres").Preload("Actors").First(&movie, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query movie"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":               movie.ID,
		"title":            movie.Title,
		"poster_base64":    movie.PosterBase64,
		"duration_minutes": movie.DurationMinutes,
		"synopsis":         movie.Synopsis,
		"release_year":     movie.ReleaseYear,
		"rating":           movie.Rating,
		"views":            movie.Views,
		"genres":           genreIDs(movie.Genres),
		"actors":           actorIDs(movie.Actors),
	})
}

// UpdateMovie - PATCH /movies/:id (auth required)
func (mc *MovieController) UpdateMovie(c *gin.Context) {
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	id := uint(id64)

	var movie models.Movie
	if err := mc.DB.Preload("Genres").Preload("Actors").First(&movie, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query movie"})
		return
	}

	var req updateMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields if provided
	if req.Title != nil {
		movie.Title = *req.Title
	}
	if req.PosterBase64 != nil {
		// if provided (even empty string), replace
		movie.PosterBase64 = *req.PosterBase64
	}
	if req.DurationMinutes != nil {
		movie.DurationMinutes = *req.DurationMinutes
	}
	if req.Synopsis != nil {
		movie.Synopsis = *req.Synopsis
	}
	if req.ReleaseYear != nil {
		movie.ReleaseYear = *req.ReleaseYear
	}
	if req.Rating != nil {
		movie.Rating = *req.Rating
	}
	if req.Views != nil {
		movie.Views = *req.Views
	}

	// handle genres replacement
	if req.Genres != nil {
		// if empty slice provided, clear relation
		if len(req.Genres) == 0 {
			if err := mc.DB.Model(&movie).Association("Genres").Clear(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to clear genres"})
				return
			}
		} else {
			var genres []models.Genre
			if err := mc.DB.Find(&genres, req.Genres).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query genres"})
				return
			}
			if len(genres) != len(req.Genres) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "one or more genres not found"})
				return
			}
			if err := mc.DB.Model(&movie).Association("Genres").Replace(&genres); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to replace genres"})
				return
			}
		}
	}

	// handle actors replacement
	if req.Actors != nil {
		if len(req.Actors) == 0 {
			if err := mc.DB.Model(&movie).Association("Actors").Clear(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to clear actors"})
				return
			}
		} else {
			var actors []models.Actor
			if err := mc.DB.Find(&actors, req.Actors).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query actors"})
				return
			}
			if len(actors) != len(req.Actors) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "one or more actors not found"})
				return
			}
			if err := mc.DB.Model(&movie).Association("Actors").Replace(&actors); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to replace actors"})
				return
			}
		}
	}

	// Save movie
	if err := mc.DB.Save(&movie).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update movie"})
		return
	}

	// reload associations
	var outGenres []models.Genre
	var outActors []models.Actor
	mc.DB.Model(&movie).Association("Genres").Find(&outGenres)
	mc.DB.Model(&movie).Association("Actors").Find(&outActors)

	c.JSON(http.StatusOK, gin.H{
		"message":          "movie updated",
		"id":               movie.ID,
		"title":            movie.Title,
		"poster_base64":    movie.PosterBase64,
		"duration_minutes": movie.DurationMinutes,
		"synopsis":         movie.Synopsis,
		"release_year":     movie.ReleaseYear,
		"rating":           movie.Rating,
		"views":            movie.Views,
		"genres":           genreIDs(outGenres),
		"actors":           actorIDs(outActors),
	})
}

// DeleteMovie - DELETE /movies/:id (auth required)
func (mc *MovieController) DeleteMovie(c *gin.Context) {
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	id := uint(id64)

	var movie models.Movie
	if err := mc.DB.First(&movie, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query movie"})
		return
	}

	// clear associations first (optional)
	if err := mc.DB.Model(&movie).Association("Genres").Clear(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to clear genres"})
		return
	}
	if err := mc.DB.Model(&movie).Association("Actors").Clear(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to clear actors"})
		return
	}

	if err := mc.DB.Delete(&movie).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete movie"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "movie deleted"})
}
