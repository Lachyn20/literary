package entity

import "github.com/google/uuid"

type ExternalLink struct {
	ID       uuid.UUID
	SiteName string
	URL      string
	Category *string
	Notes    *string
}
