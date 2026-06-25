package photoarchive

import (
	"context"

	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/domain/repository"
)

type CreatePhotoArchiveUseCase struct {
	repo repository.PhotoArchiveRepository
}

func NewCreatePhotoArchiveUseCase(repo repository.PhotoArchiveRepository) *CreatePhotoArchiveUseCase {
	return &CreatePhotoArchiveUseCase{repo: repo}
}

func (u *CreatePhotoArchiveUseCase) Execute(ctx context.Context, photo *entity.PhotoArchive) error {
	if photo.ID == uuid.Nil {
		photo.ID = uuid.New()
	}
	return u.repo.Create(ctx, photo)
}

type GetPhotoArchiveUseCase struct {
	repo repository.PhotoArchiveRepository
}

func NewGetPhotoArchiveUseCase(repo repository.PhotoArchiveRepository) *GetPhotoArchiveUseCase {
	return &GetPhotoArchiveUseCase{repo: repo}
}

func (u *GetPhotoArchiveUseCase) Execute(ctx context.Context, id uuid.UUID) (*entity.PhotoArchive, error) {
	return u.repo.GetByID(ctx, id)
}

type ListPhotoArchiveUseCase struct {
	repo repository.PhotoArchiveRepository
}

func NewListPhotoArchiveUseCase(repo repository.PhotoArchiveRepository) *ListPhotoArchiveUseCase {
	return &ListPhotoArchiveUseCase{repo: repo}
}

func (u *ListPhotoArchiveUseCase) Execute(ctx context.Context) ([]*entity.PhotoArchive, error) {
	return u.repo.List(ctx)
}

type UpdatePhotoArchiveUseCase struct {
	repo repository.PhotoArchiveRepository
}

func NewUpdatePhotoArchiveUseCase(repo repository.PhotoArchiveRepository) *UpdatePhotoArchiveUseCase {
	return &UpdatePhotoArchiveUseCase{repo: repo}
}

func (u *UpdatePhotoArchiveUseCase) Execute(ctx context.Context, photo *entity.PhotoArchive) error {
	return u.repo.Update(ctx, photo)
}

type DeletePhotoArchiveUseCase struct {
	repo repository.PhotoArchiveRepository
}

func NewDeletePhotoArchiveUseCase(repo repository.PhotoArchiveRepository) *DeletePhotoArchiveUseCase {
	return &DeletePhotoArchiveUseCase{repo: repo}
}

func (u *DeletePhotoArchiveUseCase) Execute(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}
