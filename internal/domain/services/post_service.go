package services

import (
	"context"
	"errors"
	"time"

	"github.com/elaurentium/exilium-blog-backend/internal/domain/entities"
	"github.com/elaurentium/exilium-blog-backend/internal/domain/repositories"
	"github.com/google/uuid"
)

type PostService struct {
	postRepo      repositories.PostRepository
	userRepo      repositories.UserRepository
	subRepo repositories.SubRepository
}

func NewPostService(
	postRepo repositories.PostRepository,
	userRepo repositories.UserRepository,
	subRepo repositories.SubRepository,
) *PostService {
	return &PostService{
		postRepo:      postRepo,
		userRepo:      userRepo,
		subRepo: subRepo,
	}
}

func (s *PostService) CreatePost(
	ctx context.Context,
	title string,
	content string,
	userID uuid.UUID,
	subID uuid.UUID,
) (*entities.Post, error) {
	// Verificar se o subreddit existe
	subreddit, err := s.subRepo.GetByID(ctx, subID)
	if err != nil {
		return nil, errors.New("subreddit not found")
	}

	// Verificar se o usuário existe
	_, err = s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Verificar se o subreddit é privado
	if subreddit.IsPrivate {
		// Aqui poderia ter uma lógica para verificar se o usuário é membro do subreddit
		// Por simplicidade, estamos permitindo
	}

	now := time.Now()
	post := &entities.Post{
		ID:          uuid.New(),
		Title:       title,
		Content:     content,
		UserID:      userID,
		SubID: 		 subID,
		Upvotes:     0,
		Downvotes:   0,
		IsLocked:    false,
		IsPinned:    false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = s.postRepo.Create(ctx, post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) GetPost(ctx context.Context, id uuid.UUID) (*entities.Post, error) {
	return s.postRepo.GetByID(ctx, id)
}

func (s *PostService) UpdatePost(
	ctx context.Context,
	id uuid.UUID,
	userID uuid.UUID,
	title string,
	content string,
) (*entities.Post, error) {
	post, err := s.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("post not found")
	}

	// Verificar se o usuário é o proprietário do post
	if post.UserID != userID {
		return nil, errors.New("user not authorized to update this post")
	}

	// Verificar se o post está bloqueado
	if post.IsLocked {
		return nil, errors.New("post is locked and cannot be updated")
	}

	post.Title = title
	post.Content = content
	post.UpdatedAt = time.Now()

	err = s.postRepo.Update(ctx, post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) DeletePost(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	post, err := s.postRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("post not found")
	}

	// Verificar se o usuário é o proprietário do post
	if post.UserID != userID {
		// Alternativamente, poderia verificar se o usuário é moderador ou admin
		return errors.New("user not authorized to delete this post")
	}

	return s.postRepo.Delete(ctx, id)
}

func (s *PostService) UpvotePost(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error {
	return s.postRepo.UpvotePost(ctx, postID, userID)
}

func (s *PostService) DownvotePost(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error {
	return s.postRepo.DownvotePost(ctx, postID, userID)
}

func (s *PostService) RemoveVote(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error {
	return s.postRepo.RemoveVote(ctx, postID, userID)
}

func (s *PostService) GetTrendingPosts(ctx context.Context, limit int) ([]*entities.Post, error) {
	return s.postRepo.GetTrending(ctx, limit)
}

func (s *PostService) GetPostsBySub(ctx context.Context, subID uuid.UUID, limit, offset int) ([]*entities.Post, error) {
	return s.postRepo.GetBySub(ctx, subID, limit, offset)
}

func (s *PostService) GetPostsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Post, error) {
	return s.postRepo.GetByUser(ctx, userID, limit, offset)
}