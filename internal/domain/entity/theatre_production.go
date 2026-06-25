package entity

import (
	"time"

	"github.com/google/uuid"
)

type TheatreProduction struct {
	ID           uuid.UUID
	PlayTitle    string
	TheatreName  string
	PremiereDate time.Time
	Notes        *string
	CreatedAt    time.Time
}
