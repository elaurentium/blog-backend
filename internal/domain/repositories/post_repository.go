package repositories

import (
	"context"

	"github.com/elaurentium/exilium-blog-backend/internal/domain/entities"
	"github.com/google/uuid"
)

type PostRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Post, error)
	GetBySub(ctx context.Context, subredditID uuid.UUID, limit, offset int) ([]*entities.Post, error)
	GetByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Post, error)
	Create(ctx context.Context, post *entities.Post) error
	Update(ctx context.Context, post *entities.Post) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpvotePost(ctx context.Context, postID, userID uuid.UUID) error
	DownvotePost(ctx context.Context, postID, userID uuid.UUID) error
	RemoveVote(ctx context.Context, postID, userID uuid.UUID) error
	GetTrending(ctx context.Context, limit int) ([]*entities.Post, error)
	GetCommentCount(ctx context.Context, postID uuid.UUID) (int, error)
}