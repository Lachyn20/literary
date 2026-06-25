package entity

import (
	"github.com/google/uuid"
	"time"
)

type Film struct {
	ID               uuid.UUID
	Title            string
	FilmType         FilmType
	BasedOnScenario  bool
	Director         *string
	ReleaseYear      *int
	VideoPath        *string
	CreatedAt        time.Time
}
