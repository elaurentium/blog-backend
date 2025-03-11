package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/elaurentium/exilium-blog-backend/internal/domain/entities"
	"github.com/elaurentium/exilium-blog-backend/internal/domain/repositories"
)

type CommentRepository struct {
	pool *pgxpool.Pool
}

func NewCommentRepository(pool *pgxpool.Pool) repositories.CommentRepository {
	return &CommentRepository{pool: pool}
}

func (r *CommentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Comment, error) {
	query := `
		SELECT id, content, user_id, post_id, parent_id, upvotes, downvotes, created_at, updated_at, deleted_at
		FROM comments
		WHERE id = $1 AND deleted_at IS NULL
	`

	var comment entities.Comment
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&comment.ID, &comment.Content, &comment.UserID, &comment.PostID, &comment.ParentID, &comment.Upvotes, &comment.Downvotes, &comment.CreatedAt, &comment.UpdatedAt, &comment.DeletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get comment by ID: %w", err)
	}

	return &comment, nil
}

func (r *CommentRepository) GetByPost(ctx context.Context, postID uuid.UUID, limit, offset int) ([]*entities.Comment, error) {
	query := `
		SELECT id, content, user_id, post_id, parent_id, upvotes, downvotes, created_at, updated_at
		FROM comments
		WHERE post_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments by post: %w", err)
	}
	defer rows.Close()

	var comments []*entities.Comment
	for rows.Next() {
		var comment entities.Comment
		err := rows.Scan(
			&comment.ID, &comment.Content, &comment.UserID, &comment.PostID, &comment.ParentID, &comment.Upvotes, &comment.Downvotes, &comment.CreatedAt, &comment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, &comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over comments: %w", err)
	}

	return comments, nil
}

func (r *CommentRepository) GetByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Comment, error) {
	query := `
		SELECT id, content, user_id, post_id, parent_id, upvotes, downvotes, created_at, updated_at
		FROM comments
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments by user: %w", err)
	}
	defer rows.Close()

	var comments []*entities.Comment
	for rows.Next() {
		var comment entities.Comment
		err := rows.Scan(
			&comment.ID, &comment.Content, &comment.UserID, &comment.PostID, &comment.ParentID, &comment.Upvotes, &comment.Downvotes, &comment.CreatedAt, &comment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, &comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over comments: %w", err)
	}

	return comments, nil
}

func (r *CommentRepository) GetReplies(ctx context.Context, parentID uuid.UUID, limit, offset int) ([]*entities.Comment, error) {
	query := `
		SELECT id, content, user_id, post_id, parent_id, upvotes, downvotes, created_at, updated_at
		FROM comments
		WHERE parent_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, parentID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get replies: %w", err)
	}
	defer rows.Close()

	var replies []*entities.Comment
	for rows.Next() {
		var comment entities.Comment
		err := rows.Scan(
			&comment.ID, &comment.Content, &comment.UserID, &comment.PostID, &comment.ParentID, &comment.Upvotes, &comment.Downvotes, &comment.CreatedAt, &comment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reply: %w", err)
		}
		replies = append(replies, &comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over replies: %w", err)
	}

	return replies, nil
}

func (r *CommentRepository) Create(ctx context.Context, comment *entities.Comment) error {
	query := `
		INSERT INTO comments (id, content, user_id, post_id, parent_id, upvotes, downvotes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.pool.Exec(ctx, query,
		comment.ID, comment.Content, comment.UserID, comment.PostID, comment.ParentID, comment.Upvotes, comment.Downvotes, comment.CreatedAt, comment.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}

	return nil
}

func (r *CommentRepository) Update(ctx context.Context, comment *entities.Comment) error {
	query := `
		UPDATE comments
		SET content = $2, updated_at = $3
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query,
		comment.ID, comment.Content, comment.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}

	return nil
}

func (r *CommentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE comments
		SET deleted_at = NOW()
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	return nil
}

func (r *CommentRepository) UpvoteComment(ctx context.Context, commentID, userID uuid.UUID) error {
	return r.handleVote(ctx, commentID, userID, "upvote")
}

func (r *CommentRepository) DownvoteComment(ctx context.Context, commentID, userID uuid.UUID) error {
	return r.handleVote(ctx, commentID, userID, "downvote")
}

func (r *CommentRepository) RemoveVote(ctx context.Context, commentID, userID uuid.UUID) error {
	query := `
		DELETE FROM votes
		WHERE comment_id = $1 AND user_id = $2
	`

	_, err := r.pool.Exec(ctx, query, commentID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove vote: %w", err)
	}

	return nil
}

func (r *CommentRepository) handleVote(ctx context.Context, commentID, userID uuid.UUID, voteType string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Verifica se o usuário já votou neste comentário
	var existingVote string
	query := `
		SELECT type FROM votes
		WHERE comment_id = $1 AND user_id = $2
	`
	err = tx.QueryRow(ctx, query, commentID, userID).Scan(&existingVote)
	if err != nil && err != pgx.ErrNoRows {
		return fmt.Errorf("failed to check existing vote: %w", err)
	}

	// Atualiza o voto do usuário
	if existingVote == voteType {
		// Remove o voto se for o mesmo tipo
		_, err = tx.Exec(ctx, `
			DELETE FROM votes
			WHERE comment_id = $1 AND user_id = $2
		`, commentID, userID)
	} else {
		// Insere ou atualiza o voto
		_, err = tx.Exec(ctx, `
			INSERT INTO votes (id, user_id, comment_id, type, created_at, updated_at)
			VALUES ($1, $2, $3, $4, NOW(), NOW())
			ON CONFLICT (user_id, comment_id) DO UPDATE
			SET type = $4, updated_at = NOW()
		`, uuid.New(), userID, commentID, voteType)
	}
	if err != nil {
		return fmt.Errorf("failed to handle vote: %w", err)
	}

	// Atualiza a contagem de votos no comentário
	var voteUpdate string
	if voteType == "upvote" {
		voteUpdate = "upvotes = upvotes + 1"
	} else {
		voteUpdate = "downvotes = downvotes + 1"
	}
	if existingVote == "upvote" {
		voteUpdate = "upvotes = upvotes - 1"
	} else if existingVote == "downvote" {
		voteUpdate = "downvotes = downvotes - 1"
	}

	_, err = tx.Exec(ctx, `
		UPDATE comments
		SET `+voteUpdate+`
		WHERE id = $1
	`, commentID)
	if err != nil {
		return fmt.Errorf("failed to update comment vote count: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}