package repositories

import (
	"context"

	"github.com/elaurentium/exilium-blog-backend/internal/domain/entities"
	"github.com/google/uuid"
)

type SubRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Sub, error)
	GetByName(ctx context.Context, name string) (*entities.Sub, error)
	Create(ctx context.Context, sub *entities.Sub) error
	Update(ctx context.Context, sub *entities.Sub) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*entities.Sub, error)
	GetTrending(ctx context.Context, limit int) ([]*entities.Sub, error)
}