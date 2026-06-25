package broadcast

import (
	"context"

	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/domain/repository"
)

type CreateBroadcastUseCase struct {
	repo repository.BroadcastRepository
}

func NewCreateBroadcastUseCase(repo repository.BroadcastRepository) *CreateBroadcastUseCase {
	return &CreateBroadcastUseCase{repo: repo}
}

func (u *CreateBroadcastUseCase) Execute(ctx context.Context, broadcast *entity.Broadcast) error {
	if broadcast.ID == uuid.Nil {
		broadcast.ID = uuid.New()
	}
	return u.repo.Create(ctx, broadcast)
}

type GetBroadcastUseCase struct {
	repo repository.BroadcastRepository
}

func NewGetBroadcastUseCase(repo repository.BroadcastRepository) *GetBroadcastUseCase {
	return &GetBroadcastUseCase{repo: repo}
}

func (u *GetBroadcastUseCase) Execute(ctx context.Context, id uuid.UUID) (*entity.Broadcast, error) {
	return u.repo.GetByID(ctx, id)
}

type ListBroadcastsUseCase struct {
	repo repository.BroadcastRepository
}

func NewListBroadcastsUseCase(repo repository.BroadcastRepository) *ListBroadcastsUseCase {
	return &ListBroadcastsUseCase{repo: repo}
}

func (u *ListBroadcastsUseCase) Execute(ctx context.Context) ([]*entity.Broadcast, error) {
	return u.repo.List(ctx)
}

type UpdateBroadcastUseCase struct {
	repo repository.BroadcastRepository
}

func NewUpdateBroadcastUseCase(repo repository.BroadcastRepository) *UpdateBroadcastUseCase {
	return &UpdateBroadcastUseCase{repo: repo}
}

func (u *UpdateBroadcastUseCase) Execute(ctx context.Context, broadcast *entity.Broadcast) error {
	return u.repo.Update(ctx, broadcast)
}

type DeleteBroadcastUseCase struct {
	repo repository.BroadcastRepository
}

func NewDeleteBroadcastUseCase(repo repository.BroadcastRepository) *DeleteBroadcastUseCase {
	return &DeleteBroadcastUseCase{repo: repo}
}

func (u *DeleteBroadcastUseCase) Execute(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}
