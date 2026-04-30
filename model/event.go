package model

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Location      string    `gorm:"size:200;not null" json:"location"`
	Date          time.Time `gorm:"not null" json:"date"`
	Time          string    `gorm:"size:10;not null" json:"time"`
	Description   string    `gorm:"size:1000" json:"description"`
	MaxAssistants int       `gorm:"default:0" json:"maxAssistants"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
