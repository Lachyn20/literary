package entity

import (
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID               uuid.UUID
	Title            string
	BibliographicInfo *string
	CoverImagePath   *string
	PDFPath          *string
	PageCount        *int
	PublishedYear    *int
	CreatedAt        time.Time
}

type BookPhoto struct {
	ID       uuid.UUID
	BookID   uuid.UUID
	ImagePath string
}
