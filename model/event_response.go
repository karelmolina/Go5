package model

import (
	"time"

	"github.com/google/uuid"
)

type ResponseStatus string

const (
	StatusPending  ResponseStatus = "pending"
	StatusGoing    ResponseStatus = "going"
	StatusRejected ResponseStatus = "rejected"
)

type EventResponse struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	EventID   uuid.UUID      `gorm:"type:uuid;not null;index;uniqueIndex:event_user_idx" json:"eventId"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index;uniqueIndex:event_user_idx" json:"userId"`
	Status    ResponseStatus `gorm:"type:varchar(20);not null" json:"status"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

func (EventResponse) TableName() string {
	return "event_responses"
}
