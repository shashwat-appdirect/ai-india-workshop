package repository

import (
	"context"
	"errors"
	"log"
	"os"

	"ai-india-workshop-backend/internal/models"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

type Repository struct {
	client         *firestore.Client
	subcollection string
}

func NewRepository(ctx context.Context) (*Repository, error) {
	subcollection := os.Getenv("FIRESTORE_SUBCOLLECTION_ID")
	if subcollection == "" {
		return nil, errors.New("FIRESTORE_SUBCOLLECTION_ID environment variable is required")
	}

	var app *firebase.App
	var err error

	// Check if service account path is provided (for local development)
	serviceAccountPath := os.Getenv("FIREBASE_SERVICE_ACCOUNT_PATH")
	if serviceAccountPath != "" {
		// Use service account file for local development
		opt := option.WithCredentialsFile(serviceAccountPath)
		app, err = firebase.NewApp(ctx, nil, opt)
		if err != nil {
			return nil, err
		}
	} else {
		// Use Application Default Credentials (ADC) for Cloud Run
		// Cloud Run automatically provides credentials via the attached service account
		// Project ID is required when using ADC
		projectID := os.Getenv("GCP_PROJECT_ID")
		if projectID == "" {
			// Try alternative environment variable names
			projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
			if projectID == "" {
				projectID = os.Getenv("GCLOUD_PROJECT")
			}
		}
		
		var config *firebase.Config
		if projectID != "" {
			config = &firebase.Config{
				ProjectID: projectID,
			}
		}
		
		app, err = firebase.NewApp(ctx, config)
		if err != nil {
			return nil, err
		}
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	return &Repository{
		client:         client,
		subcollection: subcollection,
	}, nil
}

func (r *Repository) getSubcollectionPath(collectionName string) *firestore.CollectionRef {
	return r.client.Collection("workshops").Doc(r.subcollection).Collection(collectionName)
}

// Attendee operations
func (r *Repository) CreateAttendee(ctx context.Context, attendee *models.Attendee) error {
	attendeesRef := r.getSubcollectionPath("attendees")
	_, _, err := attendeesRef.Add(ctx, attendee)
	return err
}

func (r *Repository) GetAllAttendees(ctx context.Context) ([]*models.Attendee, error) {
	attendeesRef := r.getSubcollectionPath("attendees")
	docs, err := attendeesRef.OrderBy("createdAt", firestore.Desc).Documents(ctx).GetAll()
	if err != nil {
		return []*models.Attendee{}, err
	}

	attendees := make([]*models.Attendee, 0)
	for _, doc := range docs {
		var attendee models.Attendee
		if err := doc.DataTo(&attendee); err != nil {
			log.Printf("Error parsing attendee: %v", err)
			continue
		}
		attendee.ID = doc.Ref.ID
		attendees = append(attendees, &attendee)
	}

	return attendees, nil
}

func (r *Repository) GetAttendeeCount(ctx context.Context) (int, error) {
	attendeesRef := r.getSubcollectionPath("attendees")
	docs, err := attendeesRef.Documents(ctx).GetAll()
	if err != nil {
		return 0, err
	}
	return len(docs), nil
}

func (r *Repository) DeleteAttendee(ctx context.Context, id string) error {
	attendeesRef := r.getSubcollectionPath("attendees")
	_, err := attendeesRef.Doc(id).Delete(ctx)
	return err
}

// Speaker operations
func (r *Repository) CreateSpeaker(ctx context.Context, speaker *models.Speaker) error {
	speakersRef := r.getSubcollectionPath("speakers")
	_, _, err := speakersRef.Add(ctx, speaker)
	return err
}

func (r *Repository) GetAllSpeakers(ctx context.Context) ([]*models.Speaker, error) {
	speakersRef := r.getSubcollectionPath("speakers")
	docs, err := speakersRef.Documents(ctx).GetAll()
	if err != nil {
		return []*models.Speaker{}, err
	}

	speakers := make([]*models.Speaker, 0)
	for _, doc := range docs {
		var speaker models.Speaker
		if err := doc.DataTo(&speaker); err != nil {
			log.Printf("Error parsing speaker: %v", err)
			continue
		}
		speaker.ID = doc.Ref.ID
		speakers = append(speakers, &speaker)
	}

	return speakers, nil
}

func (r *Repository) GetSpeaker(ctx context.Context, id string) (*models.Speaker, error) {
	speakersRef := r.getSubcollectionPath("speakers")
	doc, err := speakersRef.Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}

	var speaker models.Speaker
	if err := doc.DataTo(&speaker); err != nil {
		return nil, err
	}
	speaker.ID = doc.Ref.ID
	return &speaker, nil
}

func (r *Repository) UpdateSpeaker(ctx context.Context, id string, speaker *models.Speaker) error {
	speakersRef := r.getSubcollectionPath("speakers")
	updates := []firestore.Update{
		{Path: "name", Value: speaker.Name},
		{Path: "bio", Value: speaker.Bio},
	}
	if speaker.Avatar != "" {
		updates = append(updates, firestore.Update{Path: "avatar", Value: speaker.Avatar})
	}
	if speaker.LinkedIn != "" {
		updates = append(updates, firestore.Update{Path: "linkedin", Value: speaker.LinkedIn})
	}
	if speaker.Twitter != "" {
		updates = append(updates, firestore.Update{Path: "twitter", Value: speaker.Twitter})
	}
	_, err := speakersRef.Doc(id).Update(ctx, updates)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteSpeaker(ctx context.Context, id string) error {
	speakersRef := r.getSubcollectionPath("speakers")
	_, err := speakersRef.Doc(id).Delete(ctx)
	return err
}

// Session operations
func (r *Repository) CreateSession(ctx context.Context, session *models.Session) error {
	sessionsRef := r.getSubcollectionPath("sessions")
	_, _, err := sessionsRef.Add(ctx, session)
	return err
}

func (r *Repository) GetAllSessions(ctx context.Context) ([]*models.Session, error) {
	sessionsRef := r.getSubcollectionPath("sessions")
	docs, err := sessionsRef.Documents(ctx).GetAll()
	if err != nil {
		return []*models.Session{}, err
	}

	sessions := make([]*models.Session, 0)
	for _, doc := range docs {
		var session models.Session
		if err := doc.DataTo(&session); err != nil {
			log.Printf("Error parsing session: %v", err)
			continue
		}
		session.ID = doc.Ref.ID
		sessions = append(sessions, &session)
	}

	return sessions, nil
}

func (r *Repository) GetSession(ctx context.Context, id string) (*models.Session, error) {
	sessionsRef := r.getSubcollectionPath("sessions")
	doc, err := sessionsRef.Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}

	var session models.Session
	if err := doc.DataTo(&session); err != nil {
		return nil, err
	}
	session.ID = doc.Ref.ID
	return &session, nil
}

func (r *Repository) UpdateSession(ctx context.Context, id string, session *models.Session) error {
	sessionsRef := r.getSubcollectionPath("sessions")
	_, err := sessionsRef.Doc(id).Set(ctx, session)
	return err
}

func (r *Repository) DeleteSession(ctx context.Context, id string) error {
	sessionsRef := r.getSubcollectionPath("sessions")
	_, err := sessionsRef.Doc(id).Delete(ctx)
	return err
}

// Stats operations
func (r *Repository) GetDesignationBreakdown(ctx context.Context) ([]models.DesignationCount, error) {
	attendeesRef := r.getSubcollectionPath("attendees")
	docs, err := attendeesRef.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	designationMap := make(map[string]int)
	for _, doc := range docs {
		var attendee models.Attendee
		if err := doc.DataTo(&attendee); err != nil {
			log.Printf("Error parsing attendee for stats: %v", err)
			continue
		}
		designationMap[attendee.Designation]++
	}

	var breakdown []models.DesignationCount
	for designation, count := range designationMap {
		breakdown = append(breakdown, models.DesignationCount{
			Designation: designation,
			Count:       count,
		})
	}

	return breakdown, nil
}

