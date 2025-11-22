package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ai-india-workshop-backend/internal/models"
	"ai-india-workshop-backend/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupSpeakerTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestSpeakerHandler_GetAll(t *testing.T) {
	tests := []struct {
		name           string
		speakers       []*models.Speaker
		repoError      error
		expectedStatus int
	}{
		{
			name: "successful retrieval",
			speakers: []*models.Speaker{
				{ID: "1", Name: "Speaker 1", Bio: "Bio 1"},
				{ID: "2", Name: "Speaker 2", Bio: "Bio 2"},
			},
			repoError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "empty list",
			speakers:       []*models.Speaker{},
			repoError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "nil speakers",
			speakers:       nil,
			repoError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "repository error",
			speakers:       nil,
			repoError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repository.MockRepository)
			handler := NewSpeakerHandler(mockRepo)

			mockRepo.On("GetAllSpeakers", mock.Anything).Return(tt.speakers, tt.repoError)

			r := setupSpeakerTestRouter()
			r.GET("/speakers", handler.GetAll)

			req, _ := http.NewRequest("GET", "/speakers", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var speakers []*models.Speaker
				err := json.Unmarshal(w.Body.Bytes(), &speakers)
				require.NoError(t, err)
				if tt.speakers == nil {
					assert.Equal(t, 0, len(speakers))
				} else {
					assert.Equal(t, len(tt.speakers), len(speakers))
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSpeakerHandler_Create(t *testing.T) {
	tests := []struct {
		name           string
		speaker        models.Speaker
		repoError      error
		expectedStatus int
		expectError    bool
	}{
		{
			name: "valid creation",
			speaker: models.Speaker{
				Name:     "New Speaker",
				Bio:      "Speaker Bio",
				Avatar:   "avatar.jpg",
				LinkedIn: "linkedin.com/speaker",
				Twitter:  "@speaker",
			},
			repoError:      nil,
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name: "minimal creation",
			speaker: models.Speaker{
				Name: "New Speaker",
				Bio:  "Speaker Bio",
			},
			repoError:      nil,
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name: "empty name (no validation, will succeed)",
			speaker: models.Speaker{
				Bio: "Bio",
			},
			repoError:      nil,
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name: "repository error",
			speaker: models.Speaker{
				Name: "New Speaker",
				Bio:  "Speaker Bio",
			},
			repoError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repository.MockRepository)
			handler := NewSpeakerHandler(mockRepo)

			if !tt.expectError || tt.repoError != nil {
				mockRepo.On("CreateSpeaker", mock.Anything, mock.AnythingOfType("*models.Speaker")).Return(tt.repoError)
			}

			r := setupSpeakerTestRouter()
			r.POST("/speakers", handler.Create)

			jsonBody, _ := json.Marshal(tt.speaker)
			req, _ := http.NewRequest("POST", "/speakers", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
				var speaker models.Speaker
				err := json.Unmarshal(w.Body.Bytes(), &speaker)
				require.NoError(t, err)
				assert.Equal(t, tt.speaker.Name, speaker.Name)
				assert.Equal(t, tt.speaker.Bio, speaker.Bio)
			}

			if tt.expectedStatus == http.StatusCreated || tt.repoError != nil {
				mockRepo.AssertExpectations(t)
			}
		})
	}
}

func TestSpeakerHandler_Update(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		speaker        models.Speaker
		repoError      error
		expectedStatus int
	}{
		{
			name: "valid update",
			id:   "123",
			speaker: models.Speaker{
				Name:     "Updated Speaker",
				Bio:      "Updated Bio",
				Avatar:   "new-avatar.jpg",
				LinkedIn: "linkedin.com/updated",
				Twitter:  "@updated",
			},
			repoError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name: "update without optional fields",
			id:   "123",
			speaker: models.Speaker{
				Name: "Updated Speaker",
				Bio:  "Updated Bio",
			},
			repoError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name: "empty name (no validation, will succeed)",
			id:   "123",
			speaker: models.Speaker{
				Bio: "Bio",
			},
			repoError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name: "repository error",
			id:   "123",
			speaker: models.Speaker{
				Name: "Updated Speaker",
				Bio:  "Updated Bio",
			},
			repoError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repository.MockRepository)
			handler := NewSpeakerHandler(mockRepo)

			mockRepo.On("UpdateSpeaker", mock.Anything, tt.id, mock.AnythingOfType("*models.Speaker")).Return(tt.repoError)

			r := setupSpeakerTestRouter()
			r.PUT("/speakers/:id", handler.Update)

			jsonBody, _ := json.Marshal(tt.speaker)
			req, _ := http.NewRequest("PUT", "/speakers/"+tt.id, bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var speaker models.Speaker
				err := json.Unmarshal(w.Body.Bytes(), &speaker)
				require.NoError(t, err)
				assert.Equal(t, tt.id, speaker.ID)
				assert.Equal(t, tt.speaker.Name, speaker.Name)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSpeakerHandler_Delete(t *testing.T) {
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
			handler := NewSpeakerHandler(mockRepo)

			mockRepo.On("DeleteSpeaker", mock.Anything, tt.id).Return(tt.repoError)

			r := setupSpeakerTestRouter()
			r.DELETE("/speakers/:id", handler.Delete)

			req, _ := http.NewRequest("DELETE", "/speakers/"+tt.id, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, "Speaker deleted successfully", response["message"])
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

