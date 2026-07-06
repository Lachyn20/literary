package entity

import (
	"time"

	"github.com/google/uuid"
)

type Biography struct {
	ID        uuid.UUID
	PhotoPath *string
	UpdatedAt time.Time
	Events    []BiographyEvent
}

type BiographyEvent struct {
	ID            uuid.UUID
	BiographyID   uuid.UUID
	Year          int
	TitleTk       *string
	TitleRu       *string
	TitleEn       *string
	DescriptionTk *string
	DescriptionRu *string
	DescriptionEn *string
	SortOrder     int
	CreatedAt     time.Time
}
