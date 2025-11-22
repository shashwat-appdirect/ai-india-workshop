package handlers

import (
	"net/http"
	"os"

	"ai-india-workshop-backend/internal/repository"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	repo repository.RepositoryInterface
}

func NewAdminHandler(repo repository.RepositoryInterface) *AdminHandler {
	return &AdminHandler{repo: repo}
}

func (h *AdminHandler) Login(c *gin.Context) {
	var req struct {
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Admin password not configured"})
		return
	}

	if req.Password != adminPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	session := sessions.Default(c)
	session.Set("isAdmin", true)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *AdminHandler) GetStats(c *gin.Context) {
	breakdown, err := h.repo.GetDesignationBreakdown(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"designationBreakdown": breakdown})
}

