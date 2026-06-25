package entity

import "github.com/google/uuid"

type TranslatedByAuthor struct {
	ID                 uuid.UUID
	OriginalAuthorName string
	OriginalLanguage   string
	WorkTitle          string
	Notes              *string
}

type TranslatedIntoLanguage struct {
	ID             uuid.UUID
	LanguageName   string
	TranslatorName string
	WorkTitle      string
	Notes          *string
}
