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
