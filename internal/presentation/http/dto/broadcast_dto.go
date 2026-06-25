package dto

import "time"

type BroadcastCreateRequest struct {
    Title         string     `json:"title" validate:"required,min=1,max=255"`
    BroadcastType string     `json:"broadcast_type" validate:"required,oneof=tv radio"`
    ChannelName   string     `json:"channel_name" validate:"required,min=1,max=255"`
    BroadcastDate *time.Time `json:"broadcast_date"`
    // file handled as multipart
}

type BroadcastResponse struct {
    ID            string    `json:"id"`
    Title         string    `json:"title"`
    BroadcastType string    `json:"broadcast_type"`
    ChannelName   string    `json:"channel_name"`
    BroadcastDate time.Time `json:"broadcast_date"`
    FilePath      *string   `json:"file_path,omitempty"`
    CreatedAt     time.Time `json:"created_at"`
}
