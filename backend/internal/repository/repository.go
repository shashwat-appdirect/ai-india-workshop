package repository

import (
	"context"

	"ai-india-workshop-backend/internal/models"
)

// RepositoryInterface defines the interface for repository operations
// This allows us to mock the repository in tests
type RepositoryInterface interface {
	// Attendee operations
	CreateAttendee(ctx context.Context, attendee *models.Attendee) error
	GetAllAttendees(ctx context.Context) ([]*models.Attendee, error)
	GetAttendeeCount(ctx context.Context) (int, error)
	DeleteAttendee(ctx context.Context, id string) error

	// Speaker operations
	CreateSpeaker(ctx context.Context, speaker *models.Speaker) error
	GetAllSpeakers(ctx context.Context) ([]*models.Speaker, error)
	GetSpeaker(ctx context.Context, id string) (*models.Speaker, error)
	UpdateSpeaker(ctx context.Context, id string, speaker *models.Speaker) error
	DeleteSpeaker(ctx context.Context, id string) error

	// Session operations
	CreateSession(ctx context.Context, session *models.Session) error
	GetAllSessions(ctx context.Context) ([]*models.Session, error)
	GetSession(ctx context.Context, id string) (*models.Session, error)
	UpdateSession(ctx context.Context, id string, session *models.Session) error
	DeleteSession(ctx context.Context, id string) error

	// Stats operations
	GetDesignationBreakdown(ctx context.Context) ([]models.DesignationCount, error)
}

