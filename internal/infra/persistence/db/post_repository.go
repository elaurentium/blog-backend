package db

import (
	"context"
	"github.com/elaurentium/exilium-blog-backend/internal/domain/entities"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostRepository struct {
	pool *pgxpool.Pool
}

func NewPostRepository(pool *pgxpool.Pool) *PostRepository {
	return &PostRepository{pool: pool}
}

func (r *PostRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Post, error) {
	post := &entities.Post{}
	err := r.pool.QueryRow(ctx, "SELECT id, title, content, author_id, created_at FROM posts WHERE id = $1", id).
		Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt)
	
	return post, err
}

func (r *PostRepository) GetBySub(ctx context.Context, subredditID uuid.UUID, limit, offset int) ([]*entities.Post, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, title, content, author_id, created_at FROM posts WHERE subreddit_id = $1 LIMIT $2 OFFSET $3", subredditID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*entities.Post
	for rows.Next() {
		post := &entities.Post{}
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *PostRepository) GetByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Post, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, title, content, author_id, created_at FROM posts WHERE author_id = $1 LIMIT $2 OFFSET $3", userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*entities.Post
	for rows.Next() {
		post := &entities.Post{}
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *PostRepository) Create(ctx context.Context, post *entities.Post) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO posts (id, title, content, author_id, created_at) VALUES ($1, $2, $3, $4, $5)", post.ID, post.Title, post.Content, post.UserID, post.CreatedAt)
	return err
}

func (r *PostRepository) Update(ctx context.Context, post *entities.Post) error {
	_, err := r.pool.Exec(ctx, "UPDATE posts SET title = $1, content = $2 WHERE id = $3", post.Title, post.Content, post.ID)
	return err
}

func (r *PostRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM posts WHERE id = $1", id)
	return err
}

func (r *PostRepository) UpvotePost(ctx context.Context, postID, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO post_votes (post_id, user_id, vote_type) VALUES ($1, $2, 'up') ON CONFLICT (post_id, user_id) DO UPDATE SET vote_type = 'up'", postID, userID)
	return err
}

func (r *PostRepository) DownvotePost(ctx context.Context, postID, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO post_votes (post_id, user_id, vote_type) VALUES ($1, $2, 'down') ON CONFLICT (post_id, user_id) DO UPDATE SET vote_type = 'down'", postID, userID)
	return err
}

func (r *PostRepository) RemoveVote(ctx context.Context, postID, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM post_votes WHERE post_id = $1 AND user_id = $2", postID, userID)
	return err
}

func (r *PostRepository) GetTrending(ctx context.Context, limit int) ([]*entities.Post, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, title, content, author_id, created_at FROM posts ORDER BY score DESC LIMIT $1", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*entities.Post
	for rows.Next() {
		post := &entities.Post{}
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *PostRepository) GetCommentCount(ctx context.Context, postID uuid.UUID) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM comments WHERE post_id = $1", postID).Scan(&count)
	return count, err
}

