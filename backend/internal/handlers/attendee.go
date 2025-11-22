package handlers

import (
	"net/http"
	"time"

	"ai-india-workshop-backend/internal/models"
	"ai-india-workshop-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type AttendeeHandler struct {
	repo repository.RepositoryInterface
}

func NewAttendeeHandler(repo repository.RepositoryInterface) *AttendeeHandler {
	return &AttendeeHandler{repo: repo}
}

func (h *AttendeeHandler) Register(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Email       string `json:"email" binding:"required,email"`
		Designation string `json:"designation" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	attendee := &models.Attendee{
		Name:        req.Name,
		Email:       req.Email,
		Designation: req.Designation,
		CreatedAt:   time.Now(),
	}

	if err := h.repo.CreateAttendee(c.Request.Context(), attendee); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register attendee"})
		return
	}

	c.JSON(http.StatusCreated, attendee)
}

func (h *AttendeeHandler) GetAll(c *gin.Context) {
	attendees, err := h.repo.GetAllAttendees(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attendees"})
		return
	}

	c.JSON(http.StatusOK, attendees)
}

func (h *AttendeeHandler) GetCount(c *gin.Context) {
	count, err := h.repo.GetAttendeeCount(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

func (h *AttendeeHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.repo.DeleteAttendee(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete attendee"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendee deleted successfully"})
}

