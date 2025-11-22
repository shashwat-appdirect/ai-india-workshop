package main

import (
	"context"
	"log"
	"os"

	"ai-india-workshop-backend/internal/handlers"
	"ai-india-workshop-backend/internal/middleware"
	"ai-india-workshop-backend/internal/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables - try project root first, then current directory
	if err := godotenv.Load("../.env"); err != nil {
		if err2 := godotenv.Load(".env"); err2 != nil {
			log.Println("No .env file found, using environment variables")
		}
	}

	// Initialize Firestore repository
	ctx := context.Background()
	repo, err := repository.NewRepository(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// Serve static files (frontend) if STATIC_DIR is set (for production)
	staticDir := os.Getenv("STATIC_DIR")
	if staticDir != "" {
		// Serve static assets
		r.Static("/assets", staticDir+"/assets")
		// Serve index.html for all non-API routes (SPA routing)
		r.NoRoute(func(c *gin.Context) {
			// Don't serve index.html for API routes
			path := c.Request.URL.Path
			if len(path) >= 4 && path[:4] == "/api" {
				c.JSON(404, gin.H{"error": "Not found"})
			} else {
				c.File(staticDir + "/index.html")
			}
		})
	}

	// CORS configuration
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{frontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	// Session store
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		sessionSecret = "default-secret-change-in-production"
	}
	store := cookie.NewStore([]byte(sessionSecret))
	r.Use(sessions.Sessions("admin-session", store))

	// Initialize handlers
	attendeeHandler := handlers.NewAttendeeHandler(repo)
	speakerHandler := handlers.NewSpeakerHandler(repo)
	sessionHandler := handlers.NewSessionHandler(repo)
	adminHandler := handlers.NewAdminHandler(repo)

	// Public routes
	api := r.Group("/api")
	{
		// Attendee routes
		api.POST("/attendees", attendeeHandler.Register)
		api.GET("/attendees/count", attendeeHandler.GetCount)

		// Speaker routes
		api.GET("/speakers", speakerHandler.GetAll)

		// Session routes
		api.GET("/sessions", sessionHandler.GetAll)

		// Admin auth routes (public, must be registered here before protected routes)
		api.POST("/admin/login", adminHandler.Login)
		api.POST("/admin/logout", adminHandler.Logout)
	}

	// Protected admin routes
	admin := api.Group("/admin")
	admin.Use(middleware.RequireAdmin())
	{
		admin.GET("/stats", adminHandler.GetStats)
	}

	adminProtected := api.Group("")
	adminProtected.Use(middleware.RequireAdmin())
	{
		// Attendee admin routes
		adminProtected.GET("/attendees", attendeeHandler.GetAll)
		adminProtected.DELETE("/attendees/:id", attendeeHandler.Delete)

		// Speaker admin routes
		adminProtected.POST("/speakers", speakerHandler.Create)
		adminProtected.PUT("/speakers/:id", speakerHandler.Update)
		adminProtected.DELETE("/speakers/:id", speakerHandler.Delete)

		// Session admin routes
		adminProtected.POST("/sessions", sessionHandler.Create)
		adminProtected.PUT("/sessions/:id", sessionHandler.Update)
		adminProtected.DELETE("/sessions/:id", sessionHandler.Delete)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Debug: Print all routes
	log.Println("Registered routes:")
	for _, route := range r.Routes() {
		log.Printf("  %s %s", route.Method, route.Path)
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

