package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"ai-india-workshop-backend/internal/models"
	"ai-india-workshop-backend/internal/repository"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupAdminTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	store := cookie.NewStore([]byte("test-secret-key"))
	r.Use(sessions.Sessions("admin-session", store))
	return r
}

func TestAdminHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		password       string
		adminPassword  string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "valid password",
			password:       "correct-password",
			adminPassword:  "correct-password",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "invalid password",
			password:       "wrong-password",
			adminPassword:  "correct-password",
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name:           "missing password",
			password:       "",
			adminPassword:  "correct-password",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set admin password
			os.Setenv("ADMIN_PASSWORD", tt.adminPassword)
			defer os.Unsetenv("ADMIN_PASSWORD")

			mockRepo := new(repository.MockRepository)
			handler := NewAdminHandler(mockRepo)

			r := setupAdminTestRouter()
			r.POST("/admin/login", handler.Login)

			reqBody := map[string]string{}
			if tt.password != "" {
				reqBody["password"] = tt.password
			}
			jsonBody, _ := json.Marshal(reqBody)
			req, _ := http.NewRequest("POST", "/admin/login", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if !tt.expectError && tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.True(t, response["success"].(bool))
			}
		})
	}
}

func TestAdminHandler_Login_MissingAdminPassword(t *testing.T) {
	os.Unsetenv("ADMIN_PASSWORD")

	mockRepo := new(repository.MockRepository)
	handler := NewAdminHandler(mockRepo)

	r := setupAdminTestRouter()
	r.POST("/admin/login", handler.Login)

	reqBody := map[string]string{"password": "any-password"}
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/admin/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAdminHandler_Logout(t *testing.T) {
	mockRepo := new(repository.MockRepository)
	handler := NewAdminHandler(mockRepo)

	r := setupAdminTestRouter()
	r.POST("/admin/logout", handler.Logout)

	req, _ := http.NewRequest("POST", "/admin/logout", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Logged out successfully", response["message"])
}

func TestAdminHandler_GetStats(t *testing.T) {
	tests := []struct {
		name           string
		breakdown      []models.DesignationCount
		repoError      error
		expectedStatus int
	}{
		{
			name: "successful stats retrieval",
			breakdown: []models.DesignationCount{
				{Designation: "Engineer", Count: 5},
				{Designation: "Manager", Count: 3},
			},
			repoError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "repository error",
			breakdown:      nil,
			repoError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repository.MockRepository)
			handler := NewAdminHandler(mockRepo)

			mockRepo.On("GetDesignationBreakdown", mock.Anything).Return(tt.breakdown, tt.repoError)

			r := setupAdminTestRouter()
			r.GET("/admin/stats", handler.GetStats)

			req, _ := http.NewRequest("GET", "/admin/stats", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.NotNil(t, response["designationBreakdown"])
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

