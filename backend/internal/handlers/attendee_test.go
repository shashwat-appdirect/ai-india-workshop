package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ai-india-workshop-backend/internal/models"
	"ai-india-workshop-backend/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupAttendeeTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestAttendeeHandler_Register(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]string
		repoError      error
		expectedStatus int
		expectError    bool
	}{
		{
			name: "valid registration",
			requestBody: map[string]string{
				"name":        "John Doe",
				"email":       "john@example.com",
				"designation": "Engineer",
			},
			repoError:      nil,
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name: "missing name",
			requestBody: map[string]string{
				"email":       "john@example.com",
				"designation": "Engineer",
			},
			repoError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "missing email",
			requestBody: map[string]string{
				"name":        "John Doe",
				"designation": "Engineer",
			},
			repoError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "invalid email",
			requestBody: map[string]string{
				"name":        "John Doe",
				"email":       "invalid-email",
				"designation": "Engineer",
			},
			repoError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "missing designation",
			requestBody: map[string]string{
				"name":  "John Doe",
				"email": "john@example.com",
			},
			repoError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "repository error",
			requestBody: map[string]string{
				"name":        "John Doe",
				"email":       "john@example.com",
				"designation": "Engineer",
			},
			repoError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repository.MockRepository)
			handler := NewAttendeeHandler(mockRepo)

			if !tt.expectError || tt.repoError != nil {
				mockRepo.On("CreateAttendee", mock.Anything, mock.MatchedBy(func(attendee *models.Attendee) bool {
					return attendee.Name == tt.requestBody["name"] &&
						attendee.Email == tt.requestBody["email"] &&
						attendee.Designation == tt.requestBody["designation"]
				})).Return(tt.repoError)
			}

			r := setupAttendeeTestRouter()
			r.POST("/attendees", handler.Register)

			jsonBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/attendees", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
				var attendee models.Attendee
				err := json.Unmarshal(w.Body.Bytes(), &attendee)
				require.NoError(t, err)
				assert.Equal(t, tt.requestBody["name"], attendee.Name)
				assert.Equal(t, tt.requestBody["email"], attendee.Email)
				assert.Equal(t, tt.requestBody["designation"], attendee.Designation)
				assert.False(t, attendee.CreatedAt.IsZero())
			}

			if tt.repoError == nil && tt.expectedStatus != http.StatusBadRequest {
				mockRepo.AssertExpectations(t)
			}
		})
	}
}

func TestAttendeeHandler_GetAll(t *testing.T) {
	tests := []struct {
		name           string
		attendees      []*models.Attendee
		repoError      error
		expectedStatus int
	}{
		{
			name: "successful retrieval",
			attendees: []*models.Attendee{
				{ID: "1", Name: "John Doe", Email: "john@example.com", Designation: "Engineer", CreatedAt: time.Now()},
				{ID: "2", Name: "Jane Smith", Email: "jane@example.com", Designation: "Manager", CreatedAt: time.Now()},
			},
			repoError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "empty list",
			attendees:      []*models.Attendee{},
			repoError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "repository error",
			attendees:      nil,
			repoError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repository.MockRepository)
			handler := NewAttendeeHandler(mockRepo)

			mockRepo.On("GetAllAttendees", mock.Anything).Return(tt.attendees, tt.repoError)

			r := setupAttendeeTestRouter()
			r.GET("/attendees", handler.GetAll)

			req, _ := http.NewRequest("GET", "/attendees", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var attendees []*models.Attendee
				err := json.Unmarshal(w.Body.Bytes(), &attendees)
				require.NoError(t, err)
				assert.Equal(t, len(tt.attendees), len(attendees))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAttendeeHandler_GetCount(t *testing.T) {
	tests := []struct {
		name           string
		count          int
		repoError      error
		expectedStatus int
	}{
		{
			name:           "successful count",
			count:          10,
			repoError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "zero count",
			count:          0,
			repoError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "repository error",
			count:          0,
			repoError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repository.MockRepository)
			handler := NewAttendeeHandler(mockRepo)

			mockRepo.On("GetAttendeeCount", mock.Anything).Return(tt.count, tt.repoError)

			r := setupAttendeeTestRouter()
			r.GET("/attendees/count", handler.GetCount)

			req, _ := http.NewRequest("GET", "/attendees/count", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response map[string]int
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, tt.count, response["count"])
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAttendeeHandler_Delete(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		repoError      error
		expectedStatus int
	}{
		{
			name:           "successful deletion",
			id:             "123",
			repoError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "repository error",
			id:             "123",
			repoError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repository.MockRepository)
			handler := NewAttendeeHandler(mockRepo)

			mockRepo.On("DeleteAttendee", mock.Anything, tt.id).Return(tt.repoError)

			r := setupAttendeeTestRouter()
			r.DELETE("/attendees/:id", handler.Delete)

			req, _ := http.NewRequest("DELETE", "/attendees/"+tt.id, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, "Attendee deleted successfully", response["message"])
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

