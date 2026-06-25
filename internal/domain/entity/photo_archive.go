package entity

import (
	"time"

	"github.com/google/uuid"
)

type PhotoArchive struct {
	ID          uuid.UUID
	Title       string
	ImagePath   string
	Description *string
	TakenDate   *time.Time
	Category    PhotoCategory
	CreatedAt   time.Time
}
