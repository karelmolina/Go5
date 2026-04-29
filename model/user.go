package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Username          string     `gorm:"uniqueIndex;size:30;not null" json:"username"`
	PasswordHash      string     `gorm:"not null" json:"-"`
	Role              Role       `gorm:"type:varchar(20);default:'player'" json:"role"`
	IsApproved        bool       `gorm:"default:false" json:"isApproved"`
	ApprovedAt        *time.Time `json:"approvedAt,omitempty"`
	ApprovedBy        *uuid.UUID `json:"approvedBy,omitempty"`
	FullName          string     `gorm:"size:100" json:"fullName"`
	Nickname          string     `gorm:"size:50" json:"nickname"`
	Phone             string     `gorm:"size:20" json:"phone"`
	PhotoURL          *string    `json:"photoUrl,omitempty"`
	Positions         []Position `gorm:"type:position[]" json:"positions"`
	PreferredLanguage string     `gorm:"default:'es'" json:"preferredLanguage"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
}
