package models

import "time"

type Attendee struct {
	ID          string    `json:"id" firestore:"id"`
	Name        string    `json:"name" firestore:"name"`
	Email       string    `json:"email" firestore:"email"`
	Designation string    `json:"designation" firestore:"designation"`
	CreatedAt   time.Time `json:"createdAt" firestore:"createdAt"`
}

type Speaker struct {
	ID       string `json:"id" firestore:"id"`
	Name     string `json:"name" firestore:"name"`
	Bio      string `json:"bio" firestore:"bio"`
	Avatar   string `json:"avatar,omitempty" firestore:"avatar,omitempty"`
	LinkedIn string `json:"linkedin,omitempty" firestore:"linkedin,omitempty"`
	Twitter  string `json:"twitter,omitempty" firestore:"twitter,omitempty"`
}

type Session struct {
	ID          string   `json:"id" firestore:"id"`
	Title       string   `json:"title" firestore:"title"`
	Description string   `json:"description" firestore:"description"`
	Time        string   `json:"time" firestore:"time"`
	Speakers    []string `json:"speakers" firestore:"speakers"`
}

type SessionWithSpeakers struct {
	Session
	SpeakerDetails []Speaker `json:"speakerDetails,omitempty"`
}

type AdminStats struct {
	DesignationBreakdown []DesignationCount `json:"designationBreakdown"`
}

type DesignationCount struct {
	Designation string `json:"designation"`
	Count       int    `json:"count"`
}

