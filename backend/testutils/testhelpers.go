package testutils

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// SetupTestRouter creates a Gin router in test mode
func SetupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

// SetupTestRouterWithSessions creates a Gin router with session middleware for testing
func SetupTestRouterWithSessions() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	store := cookie.NewStore([]byte("test-secret-key"))
	r.Use(sessions.Sessions("admin-session", store))
	return r
}

// CreateTestRequest creates an HTTP test request
func CreateTestRequest(method, url string, body interface{}) *http.Request {
	var req *http.Request
	if body != nil {
		req, _ = http.NewRequest(method, url, nil)
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}
	return req
}

// PerformRequest performs an HTTP request and returns the response recorder
func PerformRequest(r http.Handler, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}


