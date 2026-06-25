package entity

import (
	"time"

	"github.com/google/uuid"
)

type PersonalLetter struct {
	ID            uuid.UUID
	Title         string
	Content       string
	LetterDate    time.Time
	ScanImagePath *string
	CreatedAt     time.Time
}
