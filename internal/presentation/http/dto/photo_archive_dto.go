package dto

import "time"

type PhotoArchiveCreateRequest struct {
    Title       string   `json:"title" validate:"required,min=1,max=255"`
    Description *string  `json:"description"`
    Category    *string  `json:"category" validate:"omitempty,oneof=archive personal"`
    TakenDate   *time.Time `json:"taken_date"`
    // image handled as multipart
}

type PhotoArchiveResponse struct {
    ID          string     `json:"id"`
    Title       string     `json:"title"`
    ImagePath   string     `json:"image_path"`
    Description *string    `json:"description,omitempty"`
    TakenDate   *time.Time `json:"taken_date,omitempty"`
    Category    string     `json:"category"`
    CreatedAt   time.Time  `json:"created_at"`
}
