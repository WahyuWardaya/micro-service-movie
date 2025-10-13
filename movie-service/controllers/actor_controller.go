package controllers

import (
	"net/http"
	"strconv"

	"movie-service/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ActorController struct {
	DB *gorm.DB
}

type createActorRequest struct {
	Name        string  `json:"name" binding:"required"`
	PhotoBase64 *string `json:"photo_base64"` // optional
}

type updateActorRequest struct {
	Name        *string `json:"name"`
	PhotoBase64 *string `json:"photo_base64"`
}

// POST /actors (auth required)
func (ac *ActorController) CreateActor(c *gin.Context) {
	var req createActorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actor := models.Actor{
		Name: req.Name,
	}
	if req.PhotoBase64 != nil {
		actor.PhotoBase64 = *req.PhotoBase64
	}

	if err := ac.DB.Create(&actor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create actor"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":           actor.ID,
		"name":         actor.Name,
		"photo_base64": actor.PhotoBase64,
	})
}

// GET /actors (public)
func (ac *ActorController) ListActors(c *gin.Context) {
	var actors []models.Actor
	if err := ac.DB.Find(&actors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list actors"})
		return
	}
	c.JSON(http.StatusOK, actors)
}

// GET /actors/:id (public)
func (ac *ActorController) GetActor(c *gin.Context) {
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var actor models.Actor
	if err := ac.DB.First(&actor, uint(id64)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "actor not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query actor"})
		return
	}
	c.JSON(http.StatusOK, actor)
}

// PATCH /actors/:id (auth required)
func (ac *ActorController) UpdateActor(c *gin.Context) {
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req updateActorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var actor models.Actor
	if err := ac.DB.First(&actor, uint(id64)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "actor not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query actor"})
		return
	}

	if req.Name != nil {
		actor.Name = *req.Name
	}
	if req.PhotoBase64 != nil {
		actor.PhotoBase64 = *req.PhotoBase64
	}

	if err := ac.DB.Save(&actor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update actor"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "actor updated", "actor": actor})
}

// DELETE /actors/:id (auth required)
func (ac *ActorController) DeleteActor(c *gin.Context) {
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := ac.DB.Delete(&models.Actor{}, uint(id64)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete actor"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "actor deleted"})
}
