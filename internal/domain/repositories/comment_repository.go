package repositories

import (
	"context"

	"github.com/elaurentium/exilium-blog-backend/internal/domain/entities"
	"github.com/google/uuid"
)

type CommentRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Comment, error)
	GetByPost(ctx context.Context, postID uuid.UUID, limit, offset int) ([]*entities.Comment, error)
	GetByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Comment, error)
	GetReplies(ctx context.Context, parentID uuid.UUID, limit, offset int) ([]*entities.Comment, error)
	Create(ctx context.Context, comment *entities.Comment) error
	Update(ctx context.Context, comment *entities.Comment) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpvoteComment(ctx context.Context, commentID, userID uuid.UUID) error
	DownvoteComment(ctx context.Context, commentID, userID uuid.UUID) error
	RemoveVote(ctx context.Context, commentID, userID uuid.UUID) error
}