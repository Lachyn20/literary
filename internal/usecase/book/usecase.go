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
}

func NewUpdateBookUseCase(repo repository.BookRepository) *UpdateBookUseCase {
	return &UpdateBookUseCase{repo: repo}
}

func (u *UpdateBookUseCase) Execute(ctx context.Context, book *entity.Book) error {
	return u.repo.Update(ctx, book)
}

type DeleteBookUseCase struct {
	repo repository.BookRepository
}

func NewDeleteBookUseCase(repo repository.BookRepository) *DeleteBookUseCase {
	return &DeleteBookUseCase{repo: repo}
}

func (u *DeleteBookUseCase) Execute(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}
