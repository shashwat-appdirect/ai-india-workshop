package handlers

import (
	"net/http"

	"ai-india-workshop-backend/internal/models"
	"ai-india-workshop-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type SpeakerHandler struct {
	repo repository.RepositoryInterface
}

func NewSpeakerHandler(repo repository.RepositoryInterface) *SpeakerHandler {
	return &SpeakerHandler{repo: repo}
}

func (h *SpeakerHandler) GetAll(c *gin.Context) {
	speakers, err := h.repo.GetAllSpeakers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch speakers"})
		return
	}

	if speakers == nil {
		speakers = []*models.Speaker{}
	}
	c.JSON(http.StatusOK, speakers)
}

func (h *SpeakerHandler) Create(c *gin.Context) {
	var speaker models.Speaker
	if err := c.ShouldBindJSON(&speaker); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.CreateSpeaker(c.Request.Context(), &speaker); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create speaker"})
		return
	}

	c.JSON(http.StatusCreated, speaker)
}

func (h *SpeakerHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var speaker models.Speaker
	if err := c.ShouldBindJSON(&speaker); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.UpdateSpeaker(c.Request.Context(), id, &speaker); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update speaker"})
		return
	}

	speaker.ID = id
	c.JSON(http.StatusOK, speaker)
}

func (h *SpeakerHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.repo.DeleteSpeaker(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete speaker"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Speaker deleted successfully"})
}

