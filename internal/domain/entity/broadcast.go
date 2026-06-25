package entity

import (
	"time"

	"github.com/google/uuid"
)

type Broadcast struct {
	ID            uuid.UUID
	Title         string
	BroadcastType BroadcastType
	ChannelName   string
	BroadcastDate time.Time
	FilePath      *string
	FileType      FileType
	CreatedAt     time.Time
}
