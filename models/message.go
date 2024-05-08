package models

import (
	"time"
)

type EmailRequest struct {
	ID           uint16     `json:"id" example:"1"`
	Text         string     `json:"text" example:"I would like to join your event."`
	EventID      string     `json:"event_id" example:"pwnrxtbi9z0v"`
	EventOwnerID string     `json:"event_owner_id" example:"pwnrxtbi9z0v"`
	RequesterID  string     `json:"requester_id" example:"pwnrxtbi9z0v"`
	Approved     *bool      `json:"approved" example:"false"`
	ApprovedAt   *time.Time `json:"approved_at,omitempty" example:"2023-11-03T10:15:30Z"`
	CreatedAt    time.Time  `json:"created_at" example:"2023-11-03T10:15:30Z"`
	UpdatedAt    time.Time  `json:"updated_at" example:"2023-11-03T10:15:30Z"`
}

type EmailRequestInput struct {
	Text    string `json:"text"`
	EventID string `json:"event_id" example:"pwnrxtbi9z0v"`
}

type EmailRequestApproveInput struct {
	Approved bool `json:"approved" example:"true"`
}
