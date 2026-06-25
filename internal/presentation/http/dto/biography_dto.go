package dto

import "time"

type BiographyResponse struct {
    ID        string     `json:"id"`
    Content   string     `json:"content"`
    PhotoPath *string    `json:"photo_path,omitempty"`
    UpdatedAt time.Time  `json:"updated_at"`
}

type BiographyUpdateRequest struct {
    Content string  `json:"content" validate:"required"`
    // photo handled as multipart
}
type BiographyCreateRequest struct {
	Content string `json:"content" validate:"required"`
	// photo handled as multipart, separately
}