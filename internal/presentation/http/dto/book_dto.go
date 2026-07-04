package dto

import "time"

type BookCreateRequest struct {
	Title             string  `json:"title" validate:"required,min=1,max=255"`
	BibliographicInfo *string `json:"bibliographic_info"`
	PageCount         *int    `json:"page_count"`
	PublishedYear     *int    `json:"published_year"`
}

type BookResponse struct {
	ID                string    `json:"id"`
	Title             string    `json:"title"`
	BibliographicInfo *string   `json:"bibliographic_info,omitempty"`
	CoverImagePath    *string   `json:"cover_image_path,omitempty"`
	PDFPath           *string   `json:"pdf_path,omitempty"`
	PageCount         *int      `json:"page_count,omitempty"`
	PublishedYear     *int      `json:"published_year,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}

type BookListResponse struct {
	Status string         `json:"status"`
	Data   []BookResponse `json:"data"`
	Total  int            `json:"total"`
}
