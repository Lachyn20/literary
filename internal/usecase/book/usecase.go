package book

import (
	"context"

	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/domain/repository"
)

type CreateBookUseCase struct {
	repo repository.BookRepository
}

func NewCreateBookUseCase(repo repository.BookRepository) *CreateBookUseCase {
	return &CreateBookUseCase{repo: repo}
}

func (u *CreateBookUseCase) Execute(ctx context.Context, book *entity.Book) error {
	if book.ID == uuid.Nil {
		book.ID = uuid.New()
	}
	return u.repo.Create(ctx, book)
}

type GetBookUseCase struct {
	repo repository.BookRepository
}

func NewGetBookUseCase(repo repository.BookRepository) *GetBookUseCase {
	return &GetBookUseCase{repo: repo}
}

func (u *GetBookUseCase) Execute(ctx context.Context, id uuid.UUID) (*entity.Book, error) {
	return u.repo.GetByID(ctx, id)
}

type ListBooksUseCase struct {
	repo repository.BookRepository
}

func NewListBooksUseCase(repo repository.BookRepository) *ListBooksUseCase {
	return &ListBooksUseCase{repo: repo}
}

func (u *ListBooksUseCase) Execute(ctx context.Context) ([]*entity.Book, error) {
	return u.repo.List(ctx)
}

type UpdateBookUseCase struct {
	repo repository.BookRepository
	fileStore repository.FileStorage
}

func NewUpdateBookUseCase(repo repository.BookRepository, fileStore repository.FileStorage) *UpdateBookUseCase {
	return &UpdateBookUseCase{repo: repo, fileStore: fileStore}
}

func (u *UpdateBookUseCase) Execute(ctx context.Context, book *entity.Book) error {
	// Kone yazgyny al, fayl yollary denesdirmek ucin
	old, err := u.repo.GetByID(ctx, book.ID)
	if err != nil {
		return err
	}
	if err := u.repo.Update(ctx, book); err != nil {
		return err
	}

	// Eger PDF fayl update bolan bolsa( old new den tapawutly bolsa ) old delete et
	if old.PDFPath != nil && (book.PDFPath == nil || *old.PDFPath != *book.PDFPath) {
		_ = u.fileStore.Remove(*old.PDFPath)
	}

	// Eger yuz ucin surat update bolan bolan bolsa suraty  delete et
	if old.CoverImagePath != nil && (book.CoverImagePath == nil || *old.CoverImagePath != *book.CoverImagePath) {
		_ = u.fileStore.Remove(*old.CoverImagePath)
	}

	return nil

}

type DeleteBookUseCase struct {
	repo repository.BookRepository
	photo repository.BookPhotoRepository
	fileStore repository.FileStorage
}

func NewDeleteBookUseCase(repo repository.BookRepository, photoRepo repository.BookPhotoRepository, fileStore repository.FileStorage) *DeleteBookUseCase {
	return &DeleteBookUseCase{repo: repo, photo: photoRepo, fileStore: fileStore}
}

func (u *DeleteBookUseCase) Execute(ctx context.Context, id uuid.UUID) error {
	book, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Degişli book_photos-lary-da al (ImagePath-lary pozmak üçin)
	photos, err := u.photo.ListByBookID(ctx, id)
	if err != nil {
		return err
	}

	if err := u.repo.Delete(ctx, id); err != nil {
		return err
	}
    
	// Diskden faýllary poz (hata bolsa-da, esasy amal üstünlikli geçen, diňe logla)
	if book.PDFPath != nil {
		_ = u.fileStore.Remove(*book.PDFPath)
	}
	if book.CoverImagePath != nil {
		_ = u.fileStore.Remove(*book.CoverImagePath)
	}
	for _, photo := range photos {
		_ = u.fileStore.Remove(photo.ImagePath)
	}

	return nil
}
