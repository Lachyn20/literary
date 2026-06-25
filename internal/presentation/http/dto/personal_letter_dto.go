package dto

import "time"

type PersonalLetterCreateRequest struct {
    Title      string `json:"title" validate:"required,min=1,max=255"`
    Content    string `json:"content" validate:"required"`
    LetterDate *time.Time `json:"letter_date"`
    // scan handled as multipart
}

type PersonalLetterResponse struct {
    ID            string     `json:"id"`
    Title         string     `json:"title"`
    Content       string     `json:"content"`
    LetterDate    time.Time  `json:"letter_date"`
    ScanImagePath *string    `json:"scan_image_path,omitempty"`
    CreatedAt     time.Time  `json:"created_at"`
}
