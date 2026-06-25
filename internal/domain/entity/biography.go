package entity

import (
	"time"

	"github.com/google/uuid"
)

type Biography struct {
	ID        uuid.UUID
	Content   string
	PhotoPath *string
	UpdatedAt time.Time
}
