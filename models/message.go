package models

import (
	"time"
)

type EmailRequest struct {
	ID           uint16     `json:"id" example:"1"`
	Text         string     `json:"text" example:"I would like to join your event."`
	EventID      string     `json:"eventId" example:"pwnrxtbi9z0v"`
	EventOwnerID string     `json:"eventOwnerId" example:"pwnrxtbi9z0v"`
	RequesterID  string     `json:"requesterId" example:"pwnrxtbi9z0v"`
	Approved     *bool      `json:"approved" example:"false"`
	ApprovedAt   *time.Time `json:"approvedAt,omitempty" example:"2023-11-03T10:15:30Z"`
	CreatedAt    time.Time  `json:"createdAt" example:"2023-11-03T10:15:30Z"`
	UpdatedAt    time.Time  `json:"updatedAt" example:"2023-11-03T10:15:30Z"`
}

type EmailRequestInput struct {
	Text    string `json:"text"`
	EventID string `json:"eventId" example:"pwnrxtbi9z0v"`
}

type EmailRequestApproveInput struct {
	Approved bool `json:"approved" example:"true"`
}

type EmailRequestApproveResponse struct {
	EmailRequest

	RequesterEmail string `json:"requesterEmail" example:"example@domain.com"`
}

type EmailRequestResponse struct {
	EmailRequest

	RequesterName  *string `json:"requesterName,omitempty" example:"John Doe"`
	RequesterEmail *string `json:"requesterEmail,omitempty" example:"email@test.com"`

	EventOwnerName  *string `json:"eventOwnerName,omitempty" example:"Owner Name"`
	EventOwnerEmail *string `json:"eventOwnerEmail,omitempty" example:"email@test.com"`
	EventName       *string `json:"eventName,omitempty" example:"Sample Event"`
	EventLocation   *string `json:"eventLocation,omitempty" example:"Central Park"`
	EventLevel      *string `json:"eventLevel,omitempty" example:"Any"`
	EventSport      *string `json:"eventSport,omitempty" example:"Basketball"`
}
