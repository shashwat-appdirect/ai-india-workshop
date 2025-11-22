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

func setupSessionTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestSessionHandler_GetAll(t *testing.T) {
	tests := []struct {
		name           string
		sessions       []*models.Session
		speakers       []*models.Speaker
		sessionsError  error
		speakersError  error
		expectedStatus int
		expectEnriched bool
	}{
		{
			name: "successful retrieval with speaker enrichment",
			sessions: []*models.Session{
				{ID: "1", Title: "Session 1", Description: "Desc 1", Time: "10:00", Speakers: []string{"sp1", "sp2"}},
			},
			speakers: []*models.Speaker{
				{ID: "sp1", Name: "Speaker 1", Bio: "Bio 1"},
				{ID: "sp2", Name: "Speaker 2", Bio: "Bio 2"},
			},
			sessionsError:  nil,
			speakersError:  nil,
			expectedStatus: http.StatusOK,
			expectEnriched: true,
		},
		{
			name: "successful retrieval without speakers",
			sessions: []*models.Session{
				{ID: "1", Title: "Session 1", Description: "Desc 1", Time: "10:00", Speakers: []string{}},
			},
			speakers:       nil,
			sessionsError:  nil,
			speakersError:  nil,
			expectedStatus: http.StatusOK,
			expectEnriched: false,
		},
		{
			name:           "empty sessions list",
			sessions:       []*models.Session{},
			speakers:       nil,
			sessionsError:  nil,
			speakersError:  nil,
			expectedStatus: http.StatusOK,
			expectEnriched: false,
		},
		{
			name:           "nil sessions",
			sessions:       nil,
			speakers:       nil,
			sessionsError:  nil,
			speakersError:  nil,
			expectedStatus: http.StatusOK,
			expectEnriched: false,
		},
		{
			name:           "sessions repository error",
			sessions:       nil,
			speakers:       nil,
			sessionsError:  assert.AnError,
			speakersError:  nil,
			expectedStatus: http.StatusInternalServerError,
			expectEnriched: false,
		},
		{
			name: "sessions success but speakers error",
			sessions: []*models.Session{
				{ID: "1", Title: "Session 1", Description: "Desc 1", Time: "10:00", Speakers: []string{"sp1"}},
			},
			speakers:       nil,
			sessionsError:  nil,
			speakersError:  assert.AnError,
			expectedStatus: http.StatusOK,
			expectEnriched: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repository.MockRepository)
			handler := NewSessionHandler(mockRepo)

			mockRepo.On("GetAllSessions", mock.Anything).Return(tt.sessions, tt.sessionsError)
			if tt.sessionsError == nil && tt.sessions != nil && len(tt.sessions) > 0 {
				mockRepo.On("GetAllSpeakers", mock.Anything).Return(tt.speakers, tt.speakersError)
			}

			r := setupSessionTestRouter()
			r.GET("/sessions", handler.GetAll)

			req, _ := http.NewRequest("GET", "/sessions", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var enrichedSessions []models.SessionWithSpeakers
				err := json.Unmarshal(w.Body.Bytes(), &enrichedSessions)
				require.NoError(t, err)

				if tt.sessions == nil || len(tt.sessions) == 0 {
					assert.Equal(t, 0, len(enrichedSessions))
				} else {
					assert.Equal(t, len(tt.sessions), len(enrichedSessions))
					if tt.expectEnriched && tt.speakers != nil {
						// Check that speaker details are enriched
						found := false
						for _, es := range enrichedSessions {
							if len(es.SpeakerDetails) > 0 {
								found = true
								break
							}
						}
						assert.True(t, found, "Expected speaker details to be enriched")
					}
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSessionHandler_Create(t *testing.T) {
	tests := []struct {
		name           string
		session        models.Session
		repoError      error
		expectedStatus int
		expectError    bool
	}{
		{
			name: "valid creation",
			session: models.Session{
				Title:       "New Session",
				Description: "Session Description",
				Time:        "14:00",
				Speakers:    []string{"sp1"},
			},
			repoError:      nil,
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name: "empty title (no validation, will succeed)",
			session: models.Session{
				Description: "Desc",
				Time:        "14:00",
			},
			repoError:      nil,
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name: "repository error",
			session: models.Session{
				Title:       "New Session",
				Description: "Session Description",
				Time:        "14:00",
				Speakers:    []string{},
			},
			repoError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repository.MockRepository)
			handler := NewSessionHandler(mockRepo)

			if !tt.expectError || tt.repoError != nil {
				mockRepo.On("CreateSession", mock.Anything, mock.AnythingOfType("*models.Session")).Return(tt.repoError)
			}

			r := setupSessionTestRouter()
			r.POST("/sessions", handler.Create)

			jsonBody, _ := json.Marshal(tt.session)
			req, _ := http.NewRequest("POST", "/sessions", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
				var session models.Session
				err := json.Unmarshal(w.Body.Bytes(), &session)
				require.NoError(t, err)
				assert.Equal(t, tt.session.Title, session.Title)
			}

			if tt.expectedStatus == http.StatusCreated || tt.repoError != nil {
				mockRepo.AssertExpectations(t)
			}
		})
	}
}

func TestSessionHandler_Update(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		session        models.Session
		repoError      error
		expectedStatus int
	}{
		{
			name: "valid update",
			id:   "123",
			session: models.Session{
				Title:       "Updated Session",
				Description: "Updated Description",
				Time:        "15:00",
				Speakers:    []string{"sp1"},
			},
			repoError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name: "empty title (no validation, will succeed)",
			id:   "123",
			session: models.Session{
				Description: "Updated Description",
				Time:        "15:00",
			},
			repoError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name: "repository error",
			id:   "123",
			session: models.Session{
				Title:       "Updated Session",
				Description: "Updated Description",
				Time:        "15:00",
			},
			repoError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repository.MockRepository)
			handler := NewSessionHandler(mockRepo)

			if tt.expectedStatus != http.StatusBadRequest {
				mockRepo.On("UpdateSession", mock.Anything, tt.id, mock.AnythingOfType("*models.Session")).Return(tt.repoError)
			}

			r := setupSessionTestRouter()
			r.PUT("/sessions/:id", handler.Update)

			jsonBody, _ := json.Marshal(tt.session)
			req, _ := http.NewRequest("PUT", "/sessions/"+tt.id, bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var session models.Session
				err := json.Unmarshal(w.Body.Bytes(), &session)
				require.NoError(t, err)
				assert.Equal(t, tt.id, session.ID)
				assert.Equal(t, tt.session.Title, session.Title)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSessionHandler_Delete(t *testing.T) {
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
			handler := NewSessionHandler(mockRepo)

			mockRepo.On("DeleteSession", mock.Anything, tt.id).Return(tt.repoError)

			r := setupSessionTestRouter()
			r.DELETE("/sessions/:id", handler.Delete)

			req, _ := http.NewRequest("DELETE", "/sessions/"+tt.id, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, "Session deleted successfully", response["message"])
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

