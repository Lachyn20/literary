package book

import (
	"context"

	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/domain/repository"
)

type BookService struct {
	repo      repository.BookRepository
	photoRepo repository.BookPhotoRepository
	fileStore repository.FileStorage
}

func NewBookService(repo repository.BookRepository, photoRepo repository.BookPhotoRepository, fileStore repository.FileStorage) *BookService {
	return &BookService{repo: repo, photoRepo: photoRepo, fileStore: fileStore}
}

func (s *BookService) Create(ctx context.Context, book *entity.Book) error {
	if book.ID == uuid.Nil {
		book.ID = uuid.New()
	}
	return s.repo.Create(ctx, book)
}

func (s *BookService) GetByID(ctx context.Context, id uuid.UUID) (*entity.Book, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *BookService) List(ctx context.Context) ([]*entity.Book, int, error) {
	books, err := s.repo.List(ctx)
	if err != nil {
		return nil, 0, err
	}
	return books, len(books), nil
}

func (s *BookService) Update(ctx context.Context, book *entity.Book) error {
	old, err := s.repo.GetByID(ctx, book.ID)
	if err != nil {
		return err
	}
	if err := s.repo.Update(ctx, book); err != nil {
		return err
	}

	// Remove old PDF if it changed
	if old.PDFPath != nil && (book.PDFPath == nil || *old.PDFPath != *book.PDFPath) {
		_ = s.fileStore.Remove(*old.PDFPath)
	}

	// Remove old cover image if it changed
	if old.CoverImagePath != nil && (book.CoverImagePath == nil || *old.CoverImagePath != *book.CoverImagePath) {
		_ = s.fileStore.Remove(*old.CoverImagePath)
	}

	return nil
}

func (s *BookService) Delete(ctx context.Context, id uuid.UUID) error {
	book, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Get book photos to delete their image files from disk
	photos, err := s.photoRepo.ListByBookID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Delete files from disk (errors are ignored; main operation succeeded)
	if book.PDFPath != nil {
		_ = s.fileStore.Remove(*book.PDFPath)
	}
	if book.CoverImagePath != nil {
		_ = s.fileStore.Remove(*book.CoverImagePath)
	}
	for _, photo := range photos {
		_ = s.fileStore.Remove(photo.ImagePath)
	}

	return nil
}
