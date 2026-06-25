package entity

import (
	"time"

	"github.com/google/uuid"
)

type CriticismArticle struct {
	ID          uuid.UUID
	Title       string
	Content     string
	Author      string
	PublishDate time.Time
}
