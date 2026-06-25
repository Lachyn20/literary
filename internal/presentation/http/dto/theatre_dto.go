package dto

import "time"

type TheatreCreateRequest struct {
    PlayTitle   string    `json:"play_title" validate:"required,min=1,max=255"`
    TheatreName string    `json:"theatre_name"`
    PremiereDate *time.Time `json:"premiere_date"`
    Notes       *string   `json:"notes"`
}

type TheatreResponse struct {
    ID           string    `json:"id"`
    PlayTitle    string    `json:"play_title"`
    TheatreName  string    `json:"theatre_name"`
    PremiereDate time.Time `json:"premiere_date"`
    Notes        *string   `json:"notes,omitempty"`
    CreatedAt    time.Time `json:"created_at"`
}
