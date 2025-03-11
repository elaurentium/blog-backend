package services

import (
	"context"
	"errors"
	"time"

	"github.com/elaurentium/exilium-blog-backend/internal/domain/entities"
	"github.com/elaurentium/exilium-blog-backend/internal/domain/repositories"
	"github.com/google/uuid"
)

type CommentService struct {
	commentRepo repositories.CommentRepository
	postRepo    repositories.PostRepository
	userRepo    repositories.UserRepository
}

func NewCommentService(
	commentRepo repositories.CommentRepository,
	postRepo repositories.PostRepository,
	userRepo repositories.UserRepository,
) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
		userRepo:    userRepo,
	}
}

func (s *CommentService) CreateComment(
	ctx context.Context,
	content string,
	userID uuid.UUID,
	postID uuid.UUID,
	parentID *uuid.UUID,
) (*entities.Comment, error) {
	// Verificar se o post existe
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	// Verificar se o post está bloqueado
	if post.IsLocked {
		return nil, errors.New("post is locked and cannot receive comments")
	}

	// Verificar se o usuário existe
	_, err = s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Verificar se o comentário pai existe, se houver
	if parentID != nil {
		_, err = s.commentRepo.GetByID(ctx, *parentID)
		if err != nil {
			return nil, errors.New("parent comment not found")
		}
	}

	now := time.Now()
	comment := &entities.Comment{
		ID:        uuid.New(),
		Content:   content,
		UserID:    userID,
		PostID:    postID,
		ParentID:  parentID,
		Upvotes:   0,
		Downvotes: 0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = s.commentRepo.Create(ctx, comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *CommentService) GetComment(ctx context.Context, id uuid.UUID) (*entities.Comment, error) {
	return s.commentRepo.GetByID(ctx, id)
}

func (s *CommentService) UpdateComment(
	ctx context.Context,
	id uuid.UUID,
	userID uuid.UUID,
	content string,
) (*entities.Comment, error) {
	comment, err := s.commentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("comment not found")
	}

	// Verificar se o usuário é o proprietário do comentário
	if comment.UserID != userID {
		return nil, errors.New("user not authorized to update this comment")
	}

	// Verificar se o post está bloqueado
	post, err := s.postRepo.GetByID(ctx, comment.PostID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	if post.IsLocked {
		return nil, errors.New("post is locked and comments cannot be updated")
	}

	comment.Content = content
	comment.UpdatedAt = time.Now()

	err = s.commentRepo.Update(ctx, comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *CommentService) DeleteComment(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	comment, err := s.commentRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("comment not found")
	}

	// Verificar se o usuário é o proprietário do comentário
	if comment.UserID != userID {
		// Alternativamente, poderia verificar se o usuário é moderador ou admin
		return errors.New("user not authorized to delete this comment")
	}

	return s.commentRepo.Delete(ctx, id)
}

func (s *CommentService) UpvoteComment(ctx context.Context, commentID uuid.UUID, userID uuid.UUID) error {
	return s.commentRepo.UpvoteComment(ctx, commentID, userID)
}

func (s *CommentService) DownvoteComment(ctx context.Context, commentID uuid.UUID, userID uuid.UUID) error {
	return s.commentRepo.DownvoteComment(ctx, commentID, userID)
}

func (s *CommentService) RemoveVote(ctx context.Context, commentID uuid.UUID, userID uuid.UUID) error {
	return s.commentRepo.RemoveVote(ctx, commentID, userID)
}

func (s *CommentService) GetCommentsByPost(ctx context.Context, postID uuid.UUID, limit, offset int) ([]*entities.Comment, error) {
	return s.commentRepo.GetByPost(ctx, postID, limit, offset)
}

func (s *CommentService) GetReplies(ctx context.Context, parentID uuid.UUID, limit, offset int) ([]*entities.Comment, error) {
	return s.commentRepo.GetReplies(ctx, parentID, limit, offset)
}