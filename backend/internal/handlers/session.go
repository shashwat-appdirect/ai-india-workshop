package handlers

import (
	"net/http"

	"ai-india-workshop-backend/internal/models"
	"ai-india-workshop-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	repo repository.RepositoryInterface
}

func NewSessionHandler(repo repository.RepositoryInterface) *SessionHandler {
	return &SessionHandler{repo: repo}
}

func (h *SessionHandler) GetAll(c *gin.Context) {
	sessions, err := h.repo.GetAllSessions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sessions"})
		return
	}

	if sessions == nil || len(sessions) == 0 {
		c.JSON(http.StatusOK, []models.SessionWithSpeakers{})
		return
	}

	// Enrich with speaker details
	speakers, err := h.repo.GetAllSpeakers(c.Request.Context())
	if err == nil && speakers != nil && len(speakers) > 0 {
		speakerMap := make(map[string]*models.Speaker)
		for _, s := range speakers {
			speakerMap[s.ID] = s
		}

		enrichedSessions := make([]models.SessionWithSpeakers, 0)
		for _, session := range sessions {
			enriched := models.SessionWithSpeakers{Session: *session}
			for _, speakerID := range session.Speakers {
				if speaker, ok := speakerMap[speakerID]; ok {
					enriched.SpeakerDetails = append(enriched.SpeakerDetails, *speaker)
				}
			}
			enrichedSessions = append(enrichedSessions, enriched)
		}
		c.JSON(http.StatusOK, enrichedSessions)
		return
	}

	// Return sessions as SessionWithSpeakers array (empty if no sessions)
	enrichedSessions := make([]models.SessionWithSpeakers, 0)
	for _, session := range sessions {
		enrichedSessions = append(enrichedSessions, models.SessionWithSpeakers{Session: *session})
	}
	c.JSON(http.StatusOK, enrichedSessions)
}

func (h *SessionHandler) Create(c *gin.Context) {
	var session models.Session
	if err := c.ShouldBindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.CreateSession(c.Request.Context(), &session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	c.JSON(http.StatusCreated, session)
}

func (h *SessionHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var session models.Session
	if err := c.ShouldBindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.UpdateSession(c.Request.Context(), id, &session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update session"})
		return
	}

	session.ID = id
	c.JSON(http.StatusOK, session)
}

func (h *SessionHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.repo.DeleteSession(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session deleted successfully"})
}

