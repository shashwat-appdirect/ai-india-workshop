package repository

import (
	"context"

	"ai-india-workshop-backend/internal/models"

	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of RepositoryInterface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateAttendee(ctx context.Context, attendee *models.Attendee) error {
	args := m.Called(ctx, attendee)
	return args.Error(0)
}

func (m *MockRepository) GetAllAttendees(ctx context.Context) ([]*models.Attendee, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Attendee), args.Error(1)
}

func (m *MockRepository) GetAttendeeCount(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Get(0).(int), args.Error(1)
}

func (m *MockRepository) DeleteAttendee(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) CreateSpeaker(ctx context.Context, speaker *models.Speaker) error {
	args := m.Called(ctx, speaker)
	return args.Error(0)
}

func (m *MockRepository) GetAllSpeakers(ctx context.Context) ([]*models.Speaker, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Speaker), args.Error(1)
}

func (m *MockRepository) GetSpeaker(ctx context.Context, id string) (*models.Speaker, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Speaker), args.Error(1)
}

func (m *MockRepository) UpdateSpeaker(ctx context.Context, id string, speaker *models.Speaker) error {
	args := m.Called(ctx, id, speaker)
	return args.Error(0)
}

func (m *MockRepository) DeleteSpeaker(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) CreateSession(ctx context.Context, session *models.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockRepository) GetAllSessions(ctx context.Context) ([]*models.Session, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Session), args.Error(1)
}

func (m *MockRepository) GetSession(ctx context.Context, id string) (*models.Session, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Session), args.Error(1)
}

func (m *MockRepository) UpdateSession(ctx context.Context, id string, session *models.Session) error {
	args := m.Called(ctx, id, session)
	return args.Error(0)
}

func (m *MockRepository) DeleteSession(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) GetDesignationBreakdown(ctx context.Context) ([]models.DesignationCount, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.DesignationCount), args.Error(1)
}


