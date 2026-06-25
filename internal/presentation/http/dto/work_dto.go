package dto

import "time"

type WorkCreateRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	CategoryID  string `json:"category_id" validate:"required,uuid4"`
	Content     string `json:"content" validate:"required"`
	Description string `json:"description"`
	AudienceType string `json:"audience_type" validate:"required,oneof=adult children"`
	PublishYear  *int   `json:"publish_year"`
}

type WorkResponse struct {
	ID string `json:"id"`
	Title string `json:"title"`
	CategoryID string `json:"category_id"`
	Content string `json:"content"`
	Description string `json:"description,omitempty"`
	AudienceType string `json:"audience_type"`
	PublishYear *int `json:"publish_year,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
