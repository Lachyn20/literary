package entity

import (
	"time"

	"github.com/google/uuid"
)

type Work struct {
	ID           uuid.UUID
	Title        string
	CategoryID   uuid.UUID
	Content      *string
	Description  *string
	AudienceType AudienceType
	PublishYear  *int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
