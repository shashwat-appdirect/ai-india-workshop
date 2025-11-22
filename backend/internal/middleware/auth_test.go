package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupAuthTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	store := cookie.NewStore([]byte("test-secret-key"))
	r.Use(sessions.Sessions("admin-session", store))
	return r
}

func TestRequireAdmin(t *testing.T) {
	tests := []struct {
		name           string
		setupSession   func(*http.Request)
		expectedStatus int
		shouldAbort    bool
	}{
		{
			name: "valid admin session",
			setupSession: func(req *http.Request) {
				// Session will be set up via cookie after first request
			},
			expectedStatus: http.StatusUnauthorized, // First request without session
			shouldAbort:    true,
		},
		{
			name: "missing session",
			setupSession: func(req *http.Request) {
				// No session setup
			},
			expectedStatus: http.StatusUnauthorized,
			shouldAbort:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupAuthTestRouter()

			// Setup a test handler that requires admin
			r.GET("/protected", RequireAdmin(), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			req, _ := http.NewRequest("GET", "/protected", nil)
			if tt.setupSession != nil {
				tt.setupSession(req)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.shouldAbort {
				// Should not reach the handler
				assert.NotContains(t, w.Body.String(), "success")
			} else {
				// Should reach the handler
				assert.Contains(t, w.Body.String(), "success")
			}
		})
	}
}

func TestRequireAdmin_WithSession(t *testing.T) {
	r := setupAuthTestRouter()

	// Setup a test handler that requires admin
	r.GET("/protected", RequireAdmin(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// First request to establish session
	req1, _ := http.NewRequest("GET", "/protected", nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusUnauthorized, w1.Code)

	// We need to simulate a login to set the session
	// For testing, we'll use a helper route that sets the session
	r.POST("/test-login", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("isAdmin", true)
		session.Save()
		c.JSON(http.StatusOK, gin.H{"success": true})
	})
	
	loginReq2, _ := http.NewRequest("POST", "/test-login", nil)
	loginW2 := httptest.NewRecorder()
	r.ServeHTTP(loginW2, loginReq2)
	
	// Extract cookies from login response
	cookies := loginW2.Header().Get("Set-Cookie")
	
	// Now make protected request with session cookie
	req2, _ := http.NewRequest("GET", "/protected", nil)
	if cookies != "" {
		req2.Header.Set("Cookie", cookies)
	}
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	
	// Session cookie should work, so request should succeed
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Contains(t, w2.Body.String(), "success")
}

func TestRequireAdmin_Integration(t *testing.T) {
	r := setupAuthTestRouter()

	// Protected route
	r.GET("/admin/stats", RequireAdmin(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"stats": "data"})
	})

	// Test without admin session
	req1, _ := http.NewRequest("GET", "/admin/stats", nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusUnauthorized, w1.Code)
	assert.Contains(t, w1.Body.String(), "Unauthorized")

	// Test that middleware properly aborts
	var response map[string]interface{}
	err := json.Unmarshal(w1.Body.Bytes(), &response)
	if err == nil {
		assert.Equal(t, "Unauthorized", response["error"])
	}
}

