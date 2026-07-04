package dto

type TranslatedByAuthorCreateRequest struct {
	OriginalAuthorName string  `json:"original_author_name" validate:"required"`
	OriginalLanguage   string  `json:"original_language" validate:"required"`
	WorkTitle          string  `json:"work_title" validate:"required"`
	Notes              *string `json:"notes"`
}

type TranslatedByAuthorResponse struct {
	ID                 string  `json:"id"`
	OriginalAuthorName string  `json:"original_author_name"`
	OriginalLanguage   string  `json:"original_language"`
	WorkTitle          string  `json:"work_title"`
	Notes              *string `json:"notes,omitempty"`
}

type TranslatedIntoLanguageCreateRequest struct {
	LanguageName   string  `json:"language_name" validate:"required"`
	TranslatorName string  `json:"translator_name" validate:"required"`
	WorkTitle      string  `json:"work_title" validate:"required"`
	Notes          *string `json:"notes"`
}

type TranslatedIntoLanguageResponse struct {
	ID             string  `json:"id"`
	LanguageName   string  `json:"language_name"`
	TranslatorName string  `json:"translator_name"`
	WorkTitle      string  `json:"work_title"`
	Notes          *string `json:"notes,omitempty"`
}

type TranslatedByAuthorListResponse struct {
	Status string                       `json:"status"`
	Data   []TranslatedByAuthorResponse `json:"data"`
	Total  int                          `json:"total"`
}

type TranslatedIntoLanguageListResponse struct {
	Status string                           `json:"status"`
	Data   []TranslatedIntoLanguageResponse `json:"data"`
	Total  int                              `json:"total"`
}
