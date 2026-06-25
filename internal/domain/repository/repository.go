package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
)

var ErrNotFound = errors.New("not found")

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	List(ctx context.Context) ([]*entity.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *entity.RefreshToken) error
	GetByToken(ctx context.Context, token string) (*entity.RefreshToken, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

type CategoryRepository interface {
	Create(ctx context.Context, category *entity.Category) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Category, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Category, error)
	List(ctx context.Context) ([]*entity.Category, error)
	Update(ctx context.Context, category *entity.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type WorkRepository interface {
	Create(ctx context.Context, work *entity.Work) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Work, error)
	List(ctx context.Context, filter WorkFilter) ([]*entity.Work, error)
	// Search performs a full-text search and returns results with total count for pagination.
	Search(ctx context.Context, filter WorkFilter) ([]*entity.Work, int, error)
	Update(ctx context.Context, work *entity.Work) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type WorkFilter struct {
	CategoryID *uuid.UUID
	AudienceType *entity.AudienceType
	PublishYear *int
	// Full-text search query
	Search *string
	// Pagination
	Page  int
	Limit int
}

type TranslatedByAuthorRepository interface {
	Create(ctx context.Context, translation *entity.TranslatedByAuthor) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.TranslatedByAuthor, error)
	List(ctx context.Context) ([]*entity.TranslatedByAuthor, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type TranslatedIntoLanguageRepository interface {
	Create(ctx context.Context, translation *entity.TranslatedIntoLanguage) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.TranslatedIntoLanguage, error)
	List(ctx context.Context) ([]*entity.TranslatedIntoLanguage, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CriticismArticleRepository interface {
	Create(ctx context.Context, article *entity.CriticismArticle) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.CriticismArticle, error)
	List(ctx context.Context) ([]*entity.CriticismArticle, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type BookRepository interface {
	Create(ctx context.Context, book *entity.Book) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Book, error)
	List(ctx context.Context) ([]*entity.Book, error)
	Update(ctx context.Context, book *entity.Book) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type BookPhotoRepository interface {
	Create(ctx context.Context, photo *entity.BookPhoto) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.BookPhoto, error)
	ListByBookID(ctx context.Context, bookID uuid.UUID) ([]*entity.BookPhoto, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type BroadcastRepository interface {
	Create(ctx context.Context, broadcast *entity.Broadcast) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Broadcast, error)
	List(ctx context.Context) ([]*entity.Broadcast, error)
	Update(ctx context.Context, broadcast *entity.Broadcast) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type TheatreProductionRepository interface {
	Create(ctx context.Context, production *entity.TheatreProduction) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.TheatreProduction, error)
	List(ctx context.Context) ([]*entity.TheatreProduction, error)
	Update(ctx context.Context, production *entity.TheatreProduction) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type FilmRepository interface {
	Create(ctx context.Context, film *entity.Film) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Film, error)
	List(ctx context.Context) ([]*entity.Film, error)
	Update(ctx context.Context, film *entity.Film) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type PhotoArchiveRepository interface {
	Create(ctx context.Context, photo *entity.PhotoArchive) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.PhotoArchive, error)
	List(ctx context.Context) ([]*entity.PhotoArchive, error)
	Update(ctx context.Context, photo *entity.PhotoArchive) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type BiographyRepository interface {
	Create(ctx context.Context, biography *entity.Biography) error
	GetLatest(ctx context.Context) (*entity.Biography, error)
	Update(ctx context.Context, biography *entity.Biography) error
}

type PersonalLetterRepository interface {
	Create(ctx context.Context, letter *entity.PersonalLetter) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.PersonalLetter, error)
	List(ctx context.Context) ([]*entity.PersonalLetter, error)
	Update(ctx context.Context, letter *entity.PersonalLetter) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ExternalLinkRepository interface {
	Create(ctx context.Context, link *entity.ExternalLink) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ExternalLink, error)
	List(ctx context.Context) ([]*entity.ExternalLink, error)
	Update(ctx context.Context, link *entity.ExternalLink) error
	Delete(ctx context.Context, id uuid.UUID) error
}
